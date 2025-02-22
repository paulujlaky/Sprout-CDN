import Routes from "../../Routes.json";

import { GlobalStorage } from "../Main";
import type { SproutAccount } from "../Misc/Structs";

import { GetURLParameter, Log, MakeRequest, RemoveURLParameter } from "../Misc/Utils";

export async function Authorize(): Promise<void> {

    const Token = GetURLParameter("Auth"); // In the case user has just logged in

    // Assume token is either here, or in cookies. Either way, we must POST to API/Auth to see

    const Response = await MakeRequest(Routes.Authorize, { Token: Token ? Token : null });

    if (Response?.JSON?.UID) {

        // We are authorized

        Log("Info", `Welcome${CheckFirstLogin() ? "" : " back"}, ${Response.JSON.Username}!`);
        
        GlobalStorage.User = Response.JSON as SproutAccount;
        
        RemoveURLParameter("Auth"); // Removes token from URL

        return;

    }

    window.location.href = "https://sprout.software/Accounts/Login?Redirect=SproutCDNNew"; // Redirects to sprout login page


}

export function CheckFirstLogin(): boolean {

    // This is written to persistent/local storage

    const WasFirst = localStorage.getItem("OOBE") == null;

    if (WasFirst) localStorage.setItem("OOBE", "true");
    
    return WasFirst;

}