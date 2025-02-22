import type { Browser } from "../Systems/Browser";

// Type Util

export type RenderedResource<T> = { HTML: string, JSON: T };

// Internal/Frontend

export enum AnimationTimes {

    Short = 150,
    Medium = 300,
    Long = 450

}

export interface Storage {

    User: SproutAccount | null;
    Browser: Browser;

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

export interface BackendFile {

    UID: string;
    Name: string;

    Size: number;

    Private: boolean;

    Authorized: string[];

    Path: string;
    
}