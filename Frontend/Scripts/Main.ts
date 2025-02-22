import { NewDirectory } from "./Misc/API";
import type { Storage } from "./Misc/Structs";
import { GetUserDirPath } from "./Misc/Utils";

import { Authorize } from "./Systems/Auth";

import { Browser } from "./Systems/Browser";

export const GlobalStorage: Storage = {

    User: null,
    Browser: new Browser()

};

await Authorize();

// If page is still executing, we are authorized

const UserHasDir = await GlobalStorage.Browser.GoTo(GetUserDirPath());

if (!UserHasDir) {

    NewDirectory(GetUserDirPath(), "", true); // Will automatically render it

}
