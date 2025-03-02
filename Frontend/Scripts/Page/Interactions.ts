import $ from 'jquery';

import { GlobalStorage } from '../Main';

import { DeleteDir, DeleteFile, MoveDir, MoveFile, NewDirectory } from '../Misc/API';
import { FetchDir, FetchFile, ReducePath } from '../Misc/Utils';

import { HideDialog, ShowContextMenu, ShowDialog, ShowFooterMessage, ToggleDragAndDropUploadIndicator, WaitForDialogResponse, type ContextMenuOption } from './Rendering';
import type { BackendDir, BackendFile } from '../Misc/Structs';

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

    HandleContextMenuInteractions();

    WatchBrandingClick();

    WatchForAccountMenu();

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

// Account Menu

function WatchForAccountMenu(): void {

    $(".DashHeaderAccount").on("click", async () => {

        const Dialog = $(".DashViewAccountDialog") as JQuery<HTMLDialogElement>;

        ShowDialog(Dialog);

        await WaitForDialogResponse(Dialog);

        HideDialog(Dialog);

    });

}

// Creation & Files

const HandleFileUploadButton = async (): Promise<void> => {

    // Prompt user to upload new file

    const FileInput = document.createElement("input");

    FileInput.type = "file";
    FileInput.accept = "*/*";
    FileInput.multiple = true;
    
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

const FileInteractionRouter = async (Event: JQuery.MouseEventBase | JQuery.TouchEventBase, FromContext: boolean = false): Promise<void> => {

    const Target = $(Event.target).closest(".InlineFile");

    if ($(Event.target).closest(".InlineFileActions").length > 0 && !FromContext) { return; } // Check that this isn't the context menu

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

// Context Menu

const ContextMenuRouter = async (Event: JQuery.MouseEventBase | JQuery.TouchEventBase | JQuery.ContextMenuEvent): Promise<void> => {

    // Either called on ContextMenu or on click. For clicks, we must check if the clicked element is .InlineFileActions
    // If a ContextMenu event, we must see if the element is an InlineFile 

    if (Event.type == "contextmenu") { Event.preventDefault(); } // Prevent default context menu
    
    const Target = $(Event.target).closest(".InlineFileActions, .InlineFile");
    const IsDir = $(Event.target).closest(".InlineFile").hasClass("Dir");

    const FetchAssociatedItem = (): BackendFile | BackendDir | null => {

        const UID = $(Event.target).closest(".InlineFile").attr("UID");

        if (!UID) { return null; }

        return FetchFile(UID) || FetchDir(UID);

    }

    const SharedOptions: ContextMenuOption[] = [

        {

            Name: "Open",
            Icon: `<ion-icon name="open-outline"></ion-icon>`,
            Action: () => { FileInteractionRouter(Event, true); } // telling router we *meant* to call it

        },
        {

            Name: "Delete",
            Icon: `<ion-icon name="trash-outline"></ion-icon>`,
            Action: async () => {
             
                const Item = FetchAssociatedItem();
                if (!Item) { return; }

                await HandleDeleteItemAction(Item); // must wait to actually delete the item

                GlobalStorage.Browser.Refresh();

            }

        },
        {

            Name: "Copy Link",
            Icon: `<ion-icon name="share-social-outline"></ion-icon>`,
            Action: () => {

                const Item = FetchAssociatedItem();
                if (!Item) { return; }

                HandleCopyLinkAction(Item);

            }

        }
    ]

    let FileOptions: ContextMenuOption[] = [

        {

            Name: "Download",
            Icon: `<ion-icon name="cloud-download-outline"></ion-icon>`,
            Action: () => {

                const Item = FetchAssociatedItem();
                if (!Item) { return; }

                HandleDownloadFileAction(Item as BackendFile);

            }

        },

        {

            Name: "Get QR Code",
            Icon: `<ion-icon name="qr-code-outline"></ion-icon>`,
            Action: () => { 
                
                const Item = FetchAssociatedItem();
                if (!Item) { return; }

                HandleQRCodeAction(Item);

            }

        },

    ];
    
    const DirOptions = SharedOptions; // Dirs get the same (base) options as files
    FileOptions = [...SharedOptions, ...FileOptions]; // Files get the base options + file-specific options

    if (Target.hasClass("InlineFileActions") && Event.type == "click") { 

        ShowContextMenu(Event, IsDir ? DirOptions : FileOptions);
        
    } else if ((Target.hasClass("InlineFile") || Target.hasClass("InlineFileActions")) && Event.type == "contextmenu") {

        ShowContextMenu(Event, IsDir ? DirOptions : FileOptions);
        
    } else if (Event.type == "contextmenu") {
        
        if (Event.target.tagName === "HTML") {
            
            // No files, so show a custom context menu

            ShowContextMenu(Event, [

                {

                    Name: "New File",
                    Icon: `<ion-icon name="document-outline"></ion-icon>`,
                    Action: HandleFileUploadButton

                },

                {

                    Name: "New Folder",
                    Icon: `<ion-icon name="folder-outline"></ion-icon>`,
                    Action: HandleNewFolderButton

                },

                {

                    Name: "Reload Content",
                    Icon: `<ion-icon name="reload-outline"></ion-icon>`,
                    Action: () => { GlobalStorage.Browser.Refresh(); }

                },

            ]); 

        }
 
    }

}

function HandleContextMenuInteractions(): void {

    RelevantElements.Document.on("contextmenu", ContextMenuRouter);
    RelevantElements.Document.on("click", ContextMenuRouter);

}

// Drag/Drop

const IsItemDrag = (Event: JQuery.DragEventBase): boolean => { return Event.originalEvent?.dataTransfer?.types.includes("uid") ?? false; }

const GetItemDragTarget = (Event: JQuery.DragEventBase): JQuery<HTMLElement> => { return $(Event.target).closest(".Dir, .AlternateDirTarget"); }
const RemoveAllItemDropTargets = (): void => { $(".Dir, .AlternateDirTarget").removeClass("DropTarget"); }
const IsTargetItem = (Event: JQuery.DragEventBase): boolean => { return GetItemDragTarget(Event).attr("UID") == Event.originalEvent?.dataTransfer?.getData("uid"); }

const HandleDragStart = (Event: JQuery.DragEventBase): void => {

    const Target = $(Event.target);

    if (Target.hasClass("InlineFile")) {

        if (!Event.originalEvent?.dataTransfer) { return; }

        Event.originalEvent.dataTransfer.dropEffect = "move";
        Event.originalEvent.dataTransfer?.setData("UID", Target.attr("UID") || "");

        Target.hasClass("Dir") ? Event.originalEvent.dataTransfer?.setData("Type", "Dir") : Event.originalEvent.dataTransfer?.setData("Type", "File");
        
    }

}

const HandleDragOver = (Event: JQuery.DragEventBase): void => {

    Event.preventDefault();

    if (IsItemDrag(Event)) {

        RemoveAllItemDropTargets();
        
        if (GetItemDragTarget(Event).length > 0 && !IsTargetItem(Event)) {
    
            GetItemDragTarget(Event).addClass("DropTarget");
            
        }
    
    } else {

        ToggleDragAndDropUploadIndicator("Show");
        
    }

}

const HandleDrops = (Event: JQuery.DropEvent): void => {

    Event.preventDefault();

    if (IsItemDrag(Event)) {
        
        HandleMoveItemDropRequest(Event);

    } else {

        HandleFileUploadDropRequest(Event);

    }

}
    

const HandleDragLeave = (Event: JQuery.DragEventBase): void => {

    if (IsItemDrag(Event)) { return; }
    
    ToggleDragAndDropUploadIndicator("Hide");

}

// Drag/Drop Heavy Lifters 

async function HandleMoveItemDropRequest(DropEvent: JQuery.DropEvent): Promise<void> {

    RemoveAllItemDropTargets();
    
    const Target = GetItemDragTarget(DropEvent);

    if (Target.length === 0) { return; }

    // Move the file

    const ItemUID = DropEvent.originalEvent?.dataTransfer?.getData("uid");
    const Type = DropEvent.originalEvent?.dataTransfer?.getData("type") as "File" | "Dir";

    const NewDirUID = Target.attr("UID");

    if (!ItemUID || ItemUID === NewDirUID) { return; } // can't move a directory into itself

    let TargetPath = FetchDir(NewDirUID || "0")?.NormalizedPath || ReducePath(GlobalStorage.Browser.Current?.Data.NormalizedPath || "");
    
    // If the target is a file, move to the parent directory. Else, move directory to other directory

    let Resp: boolean;

    if (Type === "File") {

        const FileToMove = FetchFile(ItemUID);

        if (!FileToMove) { return; }

        Resp = await MoveFile(FileToMove.NormalizedPath, TargetPath);

    } else { // should be "Dir"

        const DirToMove = FetchDir(ItemUID);

        if (!DirToMove) { return; }

        Resp = await MoveDir(DirToMove.NormalizedPath, TargetPath);

    }

    GlobalStorage.Browser.Refresh();

    if (Resp) {

        ShowFooterMessage("Success", `Moved item to <span class="DashCodeInfill">${TargetPath}</span>`, 5_000);

    } else {

        ShowFooterMessage("Error", "Failed to move item", 5_000);

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

    RelevantElements.Document.on("dragstart", HandleDragStart);
    RelevantElements.Document.on("dragleave", HandleDragLeave);

    RelevantElements.Document.on("dragover", HandleDragOver);
    RelevantElements.Document.on("drop", HandleDrops);

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

function WatchBrandingClick(): void {
    
    $(".DashHeaderBranding").on("click", () => {

        GlobalStorage.Browser.Home();

    });

}

// Utils for Context Menu

function HandleDownloadFileAction(AssociatedFile: BackendFile): void {

    const Anchor = document.createElement("a");

    Anchor.href = AssociatedFile.URL;
    Anchor.download = AssociatedFile.Name;

    Anchor.click();

}

async function HandleDeleteItemAction(AssociatedItem: BackendFile | BackendDir): Promise<void> {

    const WasSuccess = "Size" in AssociatedItem ? await DeleteFile(AssociatedItem.Path) : await DeleteDir(AssociatedItem.Path, AssociatedItem.Name);

    if (WasSuccess) {

        ShowFooterMessage("Success", `Deleted item`, 5_000);

    } else {

        ShowFooterMessage("Error", `Failed to delete item`, 5_000);

    }

}

async function HandleCopyLinkAction(AssociatedItem: BackendFile | BackendDir): Promise<void> {

    const URL = AssociatedItem.URL;

    await navigator.clipboard.writeText(URL);

    ShowFooterMessage("Success", `Copied link to clipboard`, 5_000);

}

async function HandleQRCodeAction(AssociatedItem: BackendFile | BackendDir): Promise<void> {

    alert("Coming soon");

}