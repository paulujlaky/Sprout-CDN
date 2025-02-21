import { Directory } from "../Models.ts/Dir";

export class Browser {

    // This acts as a stack which serves to keep track of which folder a user is browsing

    public Top: Directory | null;

    constructor() {

        this.Top = null;

    }

    public GoTo(Data: any): void {

        this.Top = new Directory(Data, this.Top);
        
    }

    public GoBack(): void {

        this.Top = this.Top?.Previous || null;

    }

    public GetCurrent(): any {

        return this.Top?.Data;

    }

}