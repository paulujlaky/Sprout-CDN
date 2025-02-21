import Routes from "../../Routes.json";

import { GetURLParameter, MakeRequest, RemoveURLParameter } from "../Misc";

export async function CheckForAuth(): Promise<void> {

    const Token = GetURLParameter("Auth"); // In the case user has just logged in

    // Assume token is either here, or in cookies. Either way, we must POST to API/Auth to see

    const Response = await MakeRequest(Routes.Authorize, { Token: Token ? Token : null });

    if (Response?.JSON?.UID) {

        // We are authorized

        console.log(`Hello, ${Response.JSON.Username}`);
        
        RemoveURLParameter("Auth"); // Remove token from URL

        return;

    }

    // We are not authorized

    window.location.href = "https://sprout.software/Accounts/Login?Redirect=SproutCDNNew"; // Redirects to us


}