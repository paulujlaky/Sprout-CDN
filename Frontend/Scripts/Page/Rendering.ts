import $ from "jquery";

import { GlobalStorage } from "../Main";

import type { FullDirectory } from "../Models.ts/Dir";
import { AnimationTimes, type SproutAccount } from "../Misc/Structs";

import { WriteURL } from "../Misc/Utils";
import { RemoveDialogOnEscape, SubmitDialogOnEnter } from "./Interactions";

const RelevantElements = {

    MainArea: $(".DashMain"),

    CurrentDirIndicator: $(".DashHeaderCurrentDir"),

    Navigation: {

        Backward: $(".DashFooterButton.NavigationBackward"),
        Forward: $(".DashFooterButton.NavigationForward"),

    }

}

// File List 

export function RenderDirectory(Dir: FullDirectory): void {

    RelevantElements.MainArea.html(Dir.Contents.HTML);

    WriteURL(Dir.Data.Path.replace("Store", "Dash").replace("..", ""));

    RelevantElements.CurrentDirIndicator.text(Dir.Data.Name);

}

// Loading

export function HideLoadingView(): void {

    $(".Container.DashLoading").fadeOut(AnimationTimes.Short);

}

export function ShowFooterMessage(Mode: "Loading" | "Success" | "Info" | "Error", Message: string, HideTimeout: number = -1): void {

    // assign an ID lock

    const LockID = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);

    $(".Container.DashFooterMessageSection .DashFooterMessageText").html(Message);
    $(".Container.DashFooterMessageSection").css("display", "flex").attr("NotificationLockID", LockID);

    ($(".Container.DashFooterMessageSection .DashFooterMessageLoadingIndicator"))[Mode == "Loading" ? "show" : "hide"]();

    const Icons = {

        Success: `<ion-icon name="checkmark-circle-outline"></ion-icon>`,
        Info: `<ion-icon name="information-circle-outline"></ion-icon>`,
        Error: `<ion-icon name="warning-outline"></ion-icon>`

    }

    if (Mode != "Loading") {

        $(".Container.DashFooterMessageSection .DashFooterMessageIcon").html(Icons[Mode]);
        $(".Container.DashFooterMessageSection .DashFooterMessageIcon").css("display", "flex");
        
    } else {

        $(".Container.DashFooterMessageSection .DashFooterMessageIcon").hide();

    }
    
    if (HideTimeout != -1) {

        setTimeout(() => {

            HideFooterLoadingMessage(LockID);

        }, HideTimeout)

    };

}

export function HideFooterLoadingMessage(ExpectedID?: string): void {

    // An ID is used to lock the notification, so that it doesn't hide the wrong one

    !ExpectedID || ($(".Container.DashFooterMessageSection").attr("NotificationLockID") == ExpectedID) ? $(".Container.DashFooterMessageSection").hide() : null;

}

// Auth

export function UpdateAuthElements(Account: SproutAccount, IsFirstLogin: boolean) {

    const RelevantElements = {

        UsernameEmbeds: $(".DashHeaderMessageUser"),
        EmailEmbeds: $(".DashHeaderMessageEmail"),
        
        Welcome: {

            Back: $(".DashHeaderMessageWelcomeBack"),

        }

    }

    RelevantElements.UsernameEmbeds.text(Account.Username);
    RelevantElements.EmailEmbeds.text(Account.Email);

    if (!IsFirstLogin) RelevantElements.Welcome.Back.hide();
    
}

// Dialogs

export function ShowDialog(Dialog: JQuery<HTMLDialogElement>): void {

    $(".DashDialogOverlay").show();
    Dialog.css("display", "flex");

    // Clear inputs

    Dialog.find("input").val("").trigger("focus");
    Dialog.find("Toggle").removeClass("Active");

    const CloseBtn = Dialog.find(".DashDialogButtonCancel");
    const EnterBtn = Dialog.find(".DashDialogButtonSubmit") as JQuery<HTMLButtonElement>;

    CloseBtn.one("click", () => {

        HideDialog(Dialog);

    });

    RemoveDialogOnEscape(Dialog);
    SubmitDialogOnEnter(Dialog, EnterBtn);
    
}

export function HideDialog(Dialog: JQuery<HTMLDialogElement>): void {

    Dialog.focus();
    
    $(".DashDialogOverlay").hide();
    Dialog.hide();

    Dialog.trigger("Close");
    
}

export async function WaitForDialogResponse(Dialog: JQuery<HTMLDialogElement>): Promise<boolean> {

    return new Promise((Resolve) => {

        const CloseBtn = Dialog.find(".DashDialogButtonCancel");
        const SubmitBtn = Dialog.find(".DashDialogButtonSubmit");

        CloseBtn.one("click", () => {

            Resolve(false);

        });

        SubmitBtn.one("click", () => {

            Resolve(true);
            
        });

    });

}

// Context Menu

export interface ContextMenuOption {

    Name: string,
    Icon: string,

    Action: () => void

}

export function ShowContextMenu(Event: JQuery.MouseEventBase | JQuery.TouchEventBase, Options: ContextMenuOption[]): void {

    // Get Menu
    
    const Menu = $(".DashContextMenu");

    // Position Menu

    const MenuX = Event.pageX || 0;
    const MenuY = Event.pageY || 0;

    // Check for out of bounds

    const WindowWidth = $(window).width() || 0;
    const WindowHeight = $(window).height() || 0;

    const MenuWidth = Menu.width() || 0;
    const MenuHeight = Menu.height() || 0;

    const X = (MenuX + MenuWidth + 15) > WindowWidth ? MenuX - MenuWidth : MenuX;
    const Y = (MenuY + MenuHeight + 15) > WindowHeight ? MenuY - MenuHeight : MenuY;

    // Show Menu

    Menu.css("left", X + "px");
    Menu.css("top", Y + "px");

    Menu.css("display", "flex");

    // Listen for next click (or scroll) to hide

    $(document).one("click", () => {

        Menu.hide();

    });

    $(document).one("scroll", () => {

        Menu.hide();

    });

    // Populate Menu

    Menu.html("");

    Options.forEach((Option) => {

        const OptionElement = $(`
            
            <div class="DashContextMenuItem">

                ${Option.Name}
    
                ${Option.Icon}
    
            </div>
            
        `);

        OptionElement.on("click", () => {

            Option.Action();
            Menu.hide();

        });

        Menu.append(OptionElement);

    });

}

// Navigation

export function UpdateNavigationStates(): void {

    const CanGoBack = GlobalStorage.Browser.CanGoBack;
    const CanGoForward = GlobalStorage.Browser.CanGoForward;

    RelevantElements.Navigation.Backward.attr("Disabled", CanGoBack ? null : "");
    RelevantElements.Navigation.Forward.attr("Disabled", CanGoForward ? null : "");

}

// Misc

export function ToggleDragAndDropUploadIndicator(Mode: "Show" | "Hide"): void {

    if (Mode == "Show") $(".DashUploadDropIndicator").css("display", "flex");

    else $(".DashUploadDropIndicator").hide();
    
}