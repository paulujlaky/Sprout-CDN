import { GlobalStorage } from "../Main";

import { FullDirectory, PartialDirectory } from "../Models.ts/Dir";

/**
 * This acts as a stack which serves to keep track of which folder a user is browsing.
 * Most important piece of the CDN frontend logic.
*/
export class Browser {

    public Top: PartialDirectory | null;

    constructor() {

        this.Top = null;

    }

    // Accessors

    public get Current(): PartialDirectory | null {

        return this.Top;
        
    }

    // Methods

    public async GoTo(Path: string): Promise<boolean> {

        this.Top = new PartialDirectory(Path);

        return await this.Hydrate();
        
    }

    public async GoBack(): Promise<boolean> {

        this.Top = this.Top?.Previous || null;

        return await this.Hydrate();

    }

    // Unrelated from stack

    public async Home(): Promise<boolean> {

        // Try to get the user folder

        const UserFolderPath = `/${GlobalStorage.User?.Username || "Home"}`;

        return await this.GoTo(UserFolderPath);

    }

    // Private 
    
    private async Delegate(): Promise<boolean> {

        if (!this.Top) return false;

        const FetchedTop = await this.Top.GetFull();

        this.Top = FetchedTop;

        return !!FetchedTop;

    }

    private async Hydrate(): Promise<boolean> {

        const InitResp = await this.Delegate();

        if (!InitResp) return false;

        (this.Top as FullDirectory).Display();

        return true;

    }

}