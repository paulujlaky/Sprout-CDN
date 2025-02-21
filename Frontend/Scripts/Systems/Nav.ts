import { Directory } from "../Models.ts/Dir";

export class Navigator {

    // This acts as a stack which serves to keep track of which folder a user is browsing

    public Current: Directory | null;

    constructor() {

        this.Current = null;

    }

    public GoTo(Data: any): void {

        this.Current = new Directory(Data, this.Current);
        
    }

    public GoBack(): void {

        this.Current = this.Current?.Previous || null;

    }

    public GetCurrent(): any {

        return this.Current?.Data;

    }

}