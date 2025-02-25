import $ from 'jquery';

import { GlobalStorage } from '../Main';

import { FetchDir, FetchFile, ReducePath } from '../Misc/Utils';

import { HideDialog, ShowDialog, ShowFooterMessage, ToggleDragAndDropUploadIndicator, WaitForDialogResponse } from './Rendering';

import { MoveFile, NewDirectory } from '../Misc/API';

const RelevantElements = {

    Document: $(document),

    Buttons: {

        NewFile: $(".DashFooterButton.NewFile"),
        NewFolder: $(".DashFooterButton.NewFolder"),

    },

    Inputs: {

        Toggles: $("Toggle"),
        Inputs: $("input"),

    },

    Dialogs: {

        NewFolder: $(".DashCreateDirDialog") as JQuery<HTMLDialogElement>,

    },

    Navigation: {

        Forward: $(".DashFooterButton.NavigationForward"),
        Backward: $(".DashFooterButton.NavigationBackward"),

    }

}

export function WatchPageInteractions(): void {

    RelevantElements.Buttons.NewFile.on("click", () => {

        // Prompt user to upload new file

        const FileInput = document.createElement("input");

        FileInput.type = "file";
        FileInput.accept = "*/*";
        
        FileInput.click();

        FileInput.onchange = async () => {

            const File = (FileInput.files || [])[0];

            if (!File) { return; }

            // Upload file

            await GlobalStorage.Browser.Current?.Upload(File);

        };
        
    });

    RelevantElements.Buttons.NewFolder.on("click", async () => {

        ShowDialog(RelevantElements.Dialogs.NewFolder);

        const DidRespond = await WaitForDialogResponse(RelevantElements.Dialogs.NewFolder);

        HideDialog(RelevantElements.Dialogs.NewFolder); // hide either way
        
        if (!DidRespond) { return; } // cancelled operation

        const FolderName = RelevantElements.Dialogs.NewFolder.find("#DirName").val() as string;
        const ShouldBePrivate = RelevantElements.Dialogs.NewFolder.find("#DirPrivate").hasClass("Active");
        
        if (!FolderName) { ShowFooterMessage("Error", "Folder name cannot be empty", 5_000); return; }
        if (!GlobalStorage.Browser.Current?.Data.NormalizedPath) { ShowFooterMessage("Error", "No directory selected", 5_000); return; }

        // Create folder

        const Success = await NewDirectory(FolderName, GlobalStorage.Browser.Current.Data.NormalizedPath, ShouldBePrivate);

        Success ? ShowFooterMessage("Success", `Created new folder`, 5_000) : ShowFooterMessage("Error", `Failed to create folder`, 5_000);
        
    });

    const InlineFileOrDirInteractionRouter = async (Event: JQuery.MouseEventBase | JQuery.TouchEventBase): Promise<void> => {

        const Target = $(Event.target).closest(".InlineFile");

        if ($(Event.target).closest(".InlineFileActions ").length > 0) { return; }

        const UID = Target.attr("UID");

        if (!UID) { return; }

        const RelevantFileOrDir = FetchFile(UID) || FetchDir(UID);

        if (!RelevantFileOrDir) { return; }

        const IsFile = "Size" in RelevantFileOrDir;

        if (IsFile) {

            // File

            window.open(RelevantFileOrDir.URL, "_blank");
            
        } else {

            // Directory

            const Result = await GlobalStorage.Browser.GoTo(RelevantFileOrDir.NormalizedPath);

            if (!Result) { ShowFooterMessage("Error", "Failed to navigate to directory", 5_000); }

        }
        
    };
    
    RelevantElements.Document.on("click", InlineFileOrDirInteractionRouter);
    
    // Watch for dragged in files

    RelevantElements.Document.on("dragleave", (Event) => {

        if (Event.originalEvent?.dataTransfer?.types.includes("uid")) { return; } // File move drag

        ToggleDragAndDropUploadIndicator("Hide");

    });

    RelevantElements.Document.on("drop", (Event) => {

        if (Event.originalEvent?.dataTransfer?.types.includes("uid")) { return; } // File move drag


    });

    // Watch for file move drags

    RelevantElements.Document.on("dragstart", (Event) => {

        const Target = $(Event.target);

        if (Target.hasClass("InlineFile")) {

            if (!Event.originalEvent?.dataTransfer) { return; }

            Event.originalEvent.dataTransfer.dropEffect = "move";
        
            Event.originalEvent.dataTransfer?.setData("UID", Target.attr("UID") || "");

        }

    });

    RelevantElements.Document.on("dragover", (Event) => {

        if (Event.originalEvent?.dataTransfer?.types.includes("uid")) {

            WatchForFileMoveDrags(Event);

        } else {

            ToggleDragAndDropUploadIndicator("Show");

        }

        Event.preventDefault();


    });

    RelevantElements.Document.on("drop", (Event) => {

        if (Event.originalEvent?.dataTransfer?.types.includes("uid")) {
        
            HandleFileMoveDrop(Event);

        } else {

            HandleFileUploadDrop(Event);

        }

    });

    RelevantElements.Inputs.Toggles.on("click", (Event) => {

        const Target = $(Event.target);

        const Toggle = Target.closest("Toggle");

        Toggle.toggleClass("Active");

    });

    // Navigation

    RelevantElements.Navigation.Backward.on("click", () => {

        GlobalStorage.Browser.GoBack();

    });

    RelevantElements.Navigation.Forward.on("click", () => {

        GlobalStorage.Browser.GoForward();

    });

    // Mouse back/forward buttons

    window.addEventListener("popstate", () => {

        GlobalStorage.Browser.GoTo(window.location.pathname.replace("/Dash/", ""));

    });

    // Arrow Keys

    RelevantElements.Document.on("keydown", (Event) => {

        if (Event.key === "ArrowLeft") {

            GlobalStorage.Browser.GoBack();

        } else if (Event.key === "ArrowRight") {

            GlobalStorage.Browser.GoForward();

        }

    });

}

// Misc / Util

export function RemoveDialogOnEscape(Dialog: JQuery<HTMLDialogElement>): void {

    RelevantElements.Document.one("keydown", (Event) => {

        if (Event.key === "Escape") {

            HideDialog(Dialog);

        }

    });

}

export function SubmitDialogOnEnter(Dialog: JQuery<HTMLDialogElement>, EnterButton: JQuery<HTMLButtonElement>): void {

    const Listener = (Event: JQuery.KeyDownEvent) => {

        if (Event.key === "Enter") {

            EnterButton.trigger("click");

        }

    };

    RelevantElements.Document.on("keydown", Listener);

    Dialog.on("Close", () => {

        RelevantElements.Document.off("keydown", Listener);

    });
    
}

function WatchForFileMoveDrags(DragEvent: JQuery.DragEventBase): void {

    // Check if file is one of our own

    if (!DragEvent.originalEvent?.dataTransfer?.types.includes("uid")) { return; }

    DragEvent.preventDefault();

    // Now see if we're above a valid drop target

    const Target = $(DragEvent.target); 

    $(".Dir, .AlternateDirTarget").removeClass("DropTarget"); // Reset any previous drop targets

    if (Target.closest(".Dir").length > 0 || Target.closest(".AlternateDirTarget").length > 0) {

        Target.closest(".Dir, .AlternateDirTarget").addClass("DropTarget");
        
    }
    
}

async function HandleFileMoveDrop(DropEvent: JQuery.DropEvent): Promise<void> {

    $(".Dir, .AlternateDirTarget").removeClass("DropTarget");

    const Target = $(DropEvent.target).closest(".Dir, .AlternateDirTarget");

    if (Target.closest(".Dir").length === 0 && Target.closest(".AlternateDirTarget").length === 0) { return; }

    // Move the file

    const OriginalFileUID = DropEvent.originalEvent?.dataTransfer?.getData("uid");
    const NewDirUID = Target.attr("UID");

    if (!OriginalFileUID) { return; }

    let TargetPath = FetchDir(NewDirUID || "0")?.NormalizedPath || ReducePath(GlobalStorage.Browser.Current?.Data.NormalizedPath || "");

    const FileToMove = FetchFile(OriginalFileUID);

    if (!FileToMove || !TargetPath) { return; }

    const Resp = await MoveFile(FileToMove.NormalizedPath, TargetPath);

    GlobalStorage.Browser.Refresh();

    if (Resp) {

        ShowFooterMessage("Success", `Moved file to <span class="DashCodeInfill">${TargetPath}</span>`, 5_000);

    } else {

        ShowFooterMessage("Error", "Failed to move file", 5_000);

    }
    
}

async function HandleFileUploadDrop(DropEvent: JQuery.DropEvent): Promise<void> {

    DropEvent.preventDefault();

    ToggleDragAndDropUploadIndicator("Hide");

    const Files = DropEvent.originalEvent?.dataTransfer?.files;

    if (!Files) { return; }

    for (let i = 0; i < Files.length; i++) {

        const File = Files[i];

        await GlobalStorage.Browser.Current?.Upload(File);

    }
        
}