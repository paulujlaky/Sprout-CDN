import { GlobalStorage } from "../Main";

import { AddItemToCache, CopyToClipboard } from "../Misc/Utils";

import type { BackendDir, BackendFile, RenderedResource } from "../Misc/Structs";
import { GetDirContents, NewFile } from "../Misc/API";

import { HideFooterLoadingMessage, RenderDirectory, ShowFooterMessage } from "../Page/Rendering";

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

        if (Contents.JSON.Parent == null) return null;

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

        // Add to cache

        AddItemToCache("Dirs", this.Data.UID, this.Data);
    
    }
    
    public Display(): void {

        RenderDirectory(this);

    }

    public async Upload(IncomingFile: File, ShowMsg: boolean = true): Promise<boolean> {

        // Show Loader

        if (ShowMsg) ShowFooterMessage("Loading", `Uploading ${IncomingFile.name}`);

        const File = await NewFile(this.Data.NormalizedPath, IncomingFile);
        
        if (File) {

            // Update the directory contents

            HideFooterLoadingMessage();
            
            GlobalStorage.Browser.Refresh();

            CopyToClipboard(File.URL); // Copy to clipboard the link 

            if (ShowMsg) ShowFooterMessage("Success", `Uploaded ${IncomingFile.name}`, 5_000); // Hides after 5s
            
        } else {

            // Show error message

            if (ShowMsg) ShowFooterMessage("Error", `Failed to upload ${IncomingFile.name}`, 5_000); // Hides after 5s

        }

        return !!File;

    }
}
