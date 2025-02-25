import $ from 'jquery';
import { GlobalStorage } from '../Main';
import { FetchDir, FetchFile } from '../Misc/Utils';
import { HideDialog, ShowDialog, ShowFooterMessage, WaitForDialogResponse } from './Rendering';
import { NewDirectory } from '../Misc/API';

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

        Success ? ShowFooterMessage("Success", `Created new folder ${FolderName}`, 5_000) : ShowFooterMessage("Error", `Failed to create folder ${FolderName}`, 5_000);
        
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

            console.log("File", RelevantFileOrDir);

        } else {

            // Directory

            const Result = await GlobalStorage.Browser.GoTo(RelevantFileOrDir.NormalizedPath);

            if (!Result) { ShowFooterMessage("Error", "Failed to navigate to directory", 5_000); }

        }
        
    };
    
    RelevantElements.Document.on("click", InlineFileOrDirInteractionRouter);

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

}