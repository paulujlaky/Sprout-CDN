import $ from 'jquery';

import { GlobalStorage } from '../Main';

import { MoveFile, NewDirectory } from '../Misc/API';
import { FetchDir, FetchFile, ReducePath } from '../Misc/Utils';

import { HideDialog, ShowDialog, ShowFooterMessage, ToggleDragAndDropUploadIndicator, WaitForDialogResponse } from './Rendering';

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

// General Caller

export function WatchPageInteractions(): void {

    HandleNavigation();

    HandleCreationButtonsAndFileInteractions();

    HandleDragAndDrops();

    HandleToggleClicks();

}

// Navigation

function HandleNavigation(): void {

    // Buttons

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

// Creation & Files

const HandleFileUploadButton = async (): Promise<void> => {

    // Prompt user to upload new file

    const FileInput = document.createElement("input");

    FileInput.type = "file";
    FileInput.accept = "*/*";
    
    FileInput.click();

    FileInput.onchange = async () => {

        for (let i = 0; i < (FileInput.files?.length || 0); i++) {

            const File = FileInput.files?.item(i);

            File ? await GlobalStorage.Browser.Current?.Upload(File) : null;

        }

    };

}

const HandleNewFolderButton = async (): Promise<void> => {

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
    

}

const FileInteractionRouter = async (Event: JQuery.MouseEventBase | JQuery.TouchEventBase): Promise<void> => {

    const Target = $(Event.target).closest(".InlineFile");

    if ($(Event.target).closest(".InlineFileActions").length > 0) { return; } // Check that this isn't the context menu

    const UID = Target.attr("UID");

    if (!UID) { return; }

    const RelevantFileOrDir = FetchFile(UID) || FetchDir(UID);

    if (!RelevantFileOrDir) { return; }

    if ("Size" in RelevantFileOrDir) { // File

        window.open(RelevantFileOrDir.URL, "_blank");
        
    } else {  // Directory
        
        const Result = await GlobalStorage.Browser.GoTo(RelevantFileOrDir.NormalizedPath);

        if (!Result) { ShowFooterMessage("Error", "Failed to navigate to directory", 5_000); }

    }
    
};

function HandleCreationButtonsAndFileInteractions(): void {

    RelevantElements.Buttons.NewFile.on("click", HandleFileUploadButton);

    RelevantElements.Buttons.NewFolder.on("click", HandleNewFolderButton);

    RelevantElements.Document.on("click", FileInteractionRouter);

}

// Drag/Drop

const IsFileDrag = (Event: JQuery.DragEventBase): boolean => { return Event.originalEvent?.dataTransfer?.types.includes("uid") ?? false; }
const GetFileDragTarget = (Event: JQuery.DragEventBase): JQuery<HTMLElement> => { return $(Event.target).closest(".Dir, .AlternateDirTarget"); }
const RemoveAllFileDropTargets = (): void => { $(".Dir, .AlternateDirTarget").removeClass("DropTarget"); }

const HandleFileDragStart = (Event: JQuery.DragEventBase): void => {

    const Target = $(Event.target);

    if (Target.hasClass("InlineFile")) {

        if (!Event.originalEvent?.dataTransfer) { return; }

        Event.originalEvent.dataTransfer.dropEffect = "move";
        Event.originalEvent.dataTransfer?.setData("UID", Target.attr("UID") || "");

    }

}

const HandleFileDragOver = (Event: JQuery.DragEventBase): void => {

    Event.preventDefault();

    if (IsFileDrag(Event)) {

        RemoveAllFileDropTargets();
        
        if (GetFileDragTarget(Event).length > 0) {
    
            GetFileDragTarget(Event).addClass("DropTarget");
            
        }
    
    } else {

        ToggleDragAndDropUploadIndicator("Show");
        
    }

}

const HandleFileDrops = (Event: JQuery.DropEvent): void => {

    Event.preventDefault();

    if (IsFileDrag(Event)) {
        
        HandleMoveFileDropRequest(Event);

    } else {

        HandleFileUploadDropRequest(Event);

    }

}
    

const HandleFileDragLeave = (Event: JQuery.DragEventBase): void => {

    if (IsFileDrag(Event)) { return; }
    
    ToggleDragAndDropUploadIndicator("Hide");

}

// Drag/Drop Heavy Lifters 

async function HandleMoveFileDropRequest(DropEvent: JQuery.DropEvent): Promise<void> {

    RemoveAllFileDropTargets();
    
    const Target = GetFileDragTarget(DropEvent);

    if (Target.length === 0) { return; }

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

async function HandleFileUploadDropRequest(DropEvent: JQuery.DropEvent): Promise<void> {

    DropEvent.preventDefault();

    ToggleDragAndDropUploadIndicator("Hide");

    const Files = DropEvent.originalEvent?.dataTransfer?.files;

    if (!Files) { return; }

    for (let i = 0; i < Files.length; i++) {

        const File = Files[i];

        await GlobalStorage.Browser.Current?.Upload(File);

    }
        
}

function HandleDragAndDrops(): void {

    RelevantElements.Document.on("dragstart", HandleFileDragStart);
    RelevantElements.Document.on("dragleave", HandleFileDragLeave);

    RelevantElements.Document.on("dragover", HandleFileDragOver);
    RelevantElements.Document.on("drop", HandleFileDrops);

}

// Misc

function HandleToggleClicks(): void {

    RelevantElements.Inputs.Toggles.on("click", (Event) => {

        const Target = $(Event.target);

        const Toggle = Target.closest("Toggle");

        Toggle.toggleClass("Active");

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