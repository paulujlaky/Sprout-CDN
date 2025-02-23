import type { BackendDir, BackendFile, RenderedResource } from "../Misc/Structs";
import { GetDirContents, NewFile } from "../Misc/API";

import { HideFooterLoadingMessage, RenderDirectory, ShowFooterLoadingMessage } from "../Page/Rendering";
import { GlobalStorage } from "../Main";

export class PartialDirectory {

    public Data: Partial<BackendDir>;
    public Previous: PartialDirectory | null;

    constructor(Path: string, Previous: PartialDirectory | null = null) {

        this.Data = { Path };
        this.Previous = Previous;

    }

    public async GetFull(): Promise<FullDirectory | null> {

        if (this.Data.Path == undefined) return null;

        const Contents = await GetDirContents(this.Data.Path);

        return new FullDirectory(Contents, this.Previous);

    }

}

export class FullDirectory extends PartialDirectory {

    // This also acts as a node within the Navigation stack

    public Full: boolean = true;

    public Data: BackendDir;
    public Contents: RenderedResource<{ Files: BackendFile[], Dirs: BackendDir[] }> = { HTML: "", JSON: { Files: [], Dirs: [] } };

    constructor(Contents: RenderedResource<{ Parent: BackendDir | null, Files: BackendFile[], Dirs: BackendDir[] }>, Previous: PartialDirectory | null) {

        if (!Contents.JSON.Parent) throw new Error("PartialDirectory: Creation thrown as Parent is null"); // Should never happen, just needed for TS not to complain

        super(Contents.JSON.Parent.Path, Previous);
        
        this.Data = Contents.JSON.Parent;
        this.Contents = Contents;
    
    }
    
    public Display(): void {

        RenderDirectory(this);

    }

    public async Upload(Uploaded: File): Promise<boolean> {

        // Show Loader

        ShowFooterLoadingMessage(`Uploading ${Uploaded.name}`);

        const Success = await NewFile(this.Data.Path, Uploaded);
        
        if (Success) {

            // Update the directory contents

            HideFooterLoadingMessage();
            
            GlobalStorage.Browser.Refresh();

            ShowFooterLoadingMessage(`Uploaded ${Uploaded.name}`, false, 5_000); // Hides after 5s
            
        } else {

            // Show error message

            ShowFooterLoadingMessage(`Failed to upload ${Uploaded.name}`, false, 5_000); // Hides after 5s

        }

        return Success;

    }
}
