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

// Loading View

export function HideLoadingView(): void {

    $(".Container.DashLoading").fadeOut(AnimationTimes.Short);

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