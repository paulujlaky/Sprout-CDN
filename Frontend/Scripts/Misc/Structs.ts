// Internal/Frontend

export interface Storage {

    Account: SproutAccount | null;
    Navigator: Navigator;

}

// Auth

export interface SproutAccount {

    UUID: string;

    Username: string;
    Email: string;
    Avatar?: string;

    Flags: string[];
    
} 

// Nav/Files

export interface BackendDir {

    Name: string;
    Path: string;

    Private: boolean;

    Authorized: string[];

    SubDirLength: number;
    SubFileLength: number;

}