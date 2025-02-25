import { GlobalStorage } from "../Main";

import { NewDirectory } from "../Misc/API";

import { GetUserDirPath } from "../Misc/Utils";

import { FullDirectory, PartialDirectory } from "../Models.ts/Dir";
import { UpdateNavigationStates } from "../Page/Rendering";

/**
 * This acts as a stack which serves to keep track of which folder a user is browsing.
 * Most important piece of the CDN frontend logic.
*/
export class Browser {

    public Top: FullDirectory | null;
    private InternalTop: PartialDirectory | null;

    private History: FullDirectory[] = [];

    constructor() {

        this.Top = null;
        this.InternalTop = null;

    }

    // Accessors

    public get Current(): FullDirectory | null {

        return this.Top;
        
    }

    public get CanGoBack(): boolean {

        return !!this.Top?.Previous;

    }

    public get CanGoForward(): boolean {

        return this.History.length > 0;

    }

    // Methods

    public async GoTo(Path: string): Promise<boolean> {

        this.InternalTop = new PartialDirectory(Path, this.Top);

        return await this.Hydrate();
        
    }

    public async GoBack(): Promise<boolean> {

        // Add to history so user can go back into it

        if (this.Top) this.History.push(this.Top);

        this.InternalTop = this.Top?.Previous || null;

        return await this.Hydrate();

    }

    public async GoForward(): Promise<boolean> {

        if (this.History.length == 0) return false;

        this.InternalTop = this.History.pop() || null;
        
        return await this.Hydrate();

    }

    public async Refresh(): Promise<boolean> {

        return await this.Hydrate();

    }

    // Unrelated from stack

    public async Home(): Promise<boolean> {

        let Success: boolean = await GlobalStorage.Browser.GoTo(GetUserDirPath());
    
        if (!Success) {
        
            Success = await NewDirectory(GetUserDirPath(), "", true); // Will automatically render it
    
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

        // Now, actually do the stuff we're here for

        this.Top?.Display();

        UpdateNavigationStates();

        return true;

    }

}