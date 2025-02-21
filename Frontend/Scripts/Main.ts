import type { Storage } from "./Misc/Structs";

import { Authorize } from "./Systems/Auth";

export const GlobalStorage: Storage = {

    Account: null,
    Navigator: new Navigator()

};

Authorize();