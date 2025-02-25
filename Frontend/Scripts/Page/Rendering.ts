import $ from "jquery";

import type { FullDirectory } from "../Models.ts/Dir";

import { WriteURL } from "../Misc/Utils";

import { AnimationTimes, type SproutAccount } from "../Misc/Structs";

// File List 

const RelevantElements = {

    Main: $(".DashMain"),

    CurrentDirIndicator: $(".DashHeaderCurrentDir"),

}

export function RenderDirectory(Dir: FullDirectory): void {

    RelevantElements.Main.html(Dir.Contents.HTML);

    WriteURL(Dir.Data.Path.replace("Store", "Dash"));

    RelevantElements.CurrentDirIndicator.text(Dir.Data.Name);

}

// Loading

export function HideLoadingView(): void {

    $(".Container.DashLoading").fadeOut(AnimationTimes.Short);

}

export function ShowFooterLoadingMessage(Mode: "Loading" | "Success" | "Info" | "Error", Message: string, HideTimeout: number = -1): void {

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