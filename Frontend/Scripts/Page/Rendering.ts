import $ from "jquery";

import type { FullDirectory } from "../Models.ts/Dir";

import { WriteURL } from "../Misc/Utils";

import { AnimationTimes, type SproutAccount } from "../Misc/Structs";
import { GlobalStorage } from "../Main";

// File List 

const RelevantElements = {

    Main: $(".DashMain"),

    CurrentDirIndicator: $(".DashHeaderCurrentDir"),

    Navigation: {

        Backward: $(".DashFooterButton.NavigationBackward"),
        Forward: $(".DashFooterButton.NavigationForward"),

    }

}

export function RenderDirectory(Dir: FullDirectory): void {

    RelevantElements.Main.html(Dir.Contents.HTML);

    WriteURL(Dir.Data.Path.replace("Store", "Dash").replace("..", ""));

    RelevantElements.CurrentDirIndicator.text(Dir.Data.Name);

}

// Loading

export function HideLoadingView(): void {

    $(".Container.DashLoading").fadeOut(AnimationTimes.Short);

}

export function ShowFooterMessage(Mode: "Loading" | "Success" | "Info" | "Error", Message: string, HideTimeout: number = -1): void {

    $(".Container.DashFooterMessageSection .DashFooterMessageText").text(Message);
    $(".Container.DashFooterMessageSection").css("display", "flex");

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
    
    if (HideTimeout != -1) setTimeout(HideFooterLoadingMessage, HideTimeout);

}

export function HideFooterLoadingMessage(): void {

    $(".Container.DashFooterMessageSection").hide();

}

// Auth

export function UpdateAuthElements(Account: SproutAccount, IsFirstLogin: boolean) {

    const RelevantElements = {

        UsernameEmbeds: $(".DashHeaderMessageUser"),
        
        Welcome: {

            Back: $(".DashHeaderMessageWelcomeBack"),

        }

    }

    RelevantElements.UsernameEmbeds.text(Account.Username);

    if (!IsFirstLogin) RelevantElements.Welcome.Back.hide();
    
}

// Dialogs

export function ShowDialog(Dialog: JQuery<HTMLDialogElement>): void {

    $(".DashDialogOverlay").show();
    Dialog.css("display", "flex");

    // Clear inputs

    Dialog.find("input").val("");
    Dialog.find("Toggle").removeClass("Active");

    const CloseBtn = Dialog.find(".DashDialogButtonCancel");

    CloseBtn.one("click", () => {

        HideDialog(Dialog);

    });

}

export function HideDialog(Dialog: JQuery<HTMLDialogElement>): void {

    $(".DashDialogOverlay").hide();
    Dialog.hide();
    
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

// Navigation

export function UpdateNavigationStates(): void {

    const CanGoBack = GlobalStorage.Browser.CanGoBack;
    const CanGoForward = GlobalStorage.Browser.CanGoForward;

    RelevantElements.Navigation.Backward.attr("Disabled", CanGoBack ? null : "");
    RelevantElements.Navigation.Forward.attr("Disabled", CanGoForward ? null : "");

}