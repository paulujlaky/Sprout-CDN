import type { BackendDir } from "../Misc/Structs";

export class Directory {

    // This also acts as a node within the Navigation stack

    public Data: BackendDir;
    public Previous: Directory | null;

    constructor(Data: BackendDir, Previous: Directory | null) {

        this.Previous = Previous;
        this.Data = Data;
    
    }

}
