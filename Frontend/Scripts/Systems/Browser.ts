import { GlobalStorage } from "../Main";

import { NewDirectory } from "../Misc/API";
import { GetUserDir } from "../Misc/Utils";

import { FullDirectory, PartialDirectory } from "../Models.ts/Dir";

/**
 * This acts as a stack which serves to keep track of which folder a user is browsing.
 * Most important piece of the CDN frontend logic.
*/
export class Browser {

    public Top: FullDirectory | null;
    private InternalTop: PartialDirectory | null;

    constructor() {

        this.Top = null;
        this.InternalTop = null;

    }

    // Accessors

    public get Current(): FullDirectory | null {

        return this.Top;
        
    }

    // Methods

    public async GoTo(Path: string): Promise<boolean> {

        this.InternalTop = new PartialDirectory(Path);

        return await this.Hydrate();
        
    }

    public async GoBack(): Promise<boolean> {

        this.InternalTop = this.Top?.Previous || null;

        return await this.Hydrate();

    }

    public async Refresh(): Promise<boolean> {

        return await this.Hydrate();

    }

    // Unrelated from stack

    public async Home(): Promise<boolean> {

        let Success: boolean = await GlobalStorage.Browser.GoTo(GetUserDir());
    
        if (!Success) {
        
            Success = await NewDirectory(GetUserDir(), "", true); // Will automatically render it
    
        }
    
        return Success;

    }

    // Private 
    
    private async Delegate(): Promise<boolean> {

        if (!this.InternalTop) return false;

        const FetchedTop = await this.InternalTop.GetFull();

        this.Top = FetchedTop;

        return !!FetchedTop;

    }

    private async Hydrate(): Promise<boolean> {

        const InitResp = await this.Delegate();

        if (!InitResp) return false;

        this.Top?.Display();

        return true;

    }

}