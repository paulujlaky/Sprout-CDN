import type { Storage } from "./Misc/Structs";

import { Browser } from "./Systems/Browser";
import { Authorize } from "./Systems/Auth";

import { Log } from "./Misc/Utils";
import { HideLoadingView } from "./Page/Rendering";
import { WatchPageInteractions } from "./Page/Interactions";

export const GlobalStorage: Storage = {

    User: null,
    Browser: new Browser(),

    Cache: {

        Files: {},
        Dirs: {}

    }

};

await Authorize();

// If page is still executing, we are authorized

// We must see if we have an implicit path to go to

if (window.location.pathname == "/") {

    Log("Info", "Navigating to user home...");

    await GlobalStorage.Browser.Home();

} else {

    Log("Info", "Navigating to implicit path...");

    const WasSuccessful = await GlobalStorage.Browser.GoTo(window.location.pathname.replace("/Dash/", ""));

    if (!WasSuccessful) {

        Log("Error", "Implicit path navigation failed, navigating to user home...");

        await GlobalStorage.Browser.Home();

    }

}

Log("Info", "Initial navigation complete.");

HideLoadingView();

WatchPageInteractions();