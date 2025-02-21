import type { Storage } from "./Misc/Structs";

import { Authorize } from "./Systems/Auth";

import { Browser } from "./Systems/Nav";

export const GlobalStorage: Storage = {

    Account: null,
    Navigator: new Browser()

};

Authorize();