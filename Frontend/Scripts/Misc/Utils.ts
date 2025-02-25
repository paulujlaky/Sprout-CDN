// Window / URL

import { GlobalStorage } from "../Main";
import type { BackendDir, BackendFile } from "./Structs";

export function GetURLParameter(Name: string): string | null {

    const URLInstance = new URL(window.location.href);
    return URLInstance.searchParams.get(Name);

}

export function RemoveURLParameter(Name: string): void {

    const URLInstance = new URL(window.location.href);
    URLInstance.searchParams.delete(Name);

    WriteURL(URLInstance.toString());
    
}

export function WriteURL(URL: string): void {

    // Overwrites everything after domain

    window.history.pushState({}, "", URL);
    
}

// Network

interface BackendResponse {

    Message: string;

    JSON: any;
    HTML: string;

}

export async function MakeRequest(Route: { URL: string, Method: string }, Body?: any, Headers: { [key: string]: string } = {}, Params: { [key: string]: string } = {}, IsFormData: boolean = false): Promise<BackendResponse | null> {

    let RouteURL = Route.URL;

    const Options: RequestInit = {

        method: Route.Method,

        headers: IsFormData ? Headers : {

            "Content-Type": "application/json",
            ...Headers

        },

        body: Body ? (IsFormData ? Body : JSON.stringify(Body)) : undefined

    };

    for (const Param in Params) {

        RouteURL = Route.URL.replace(`:${Param}`, Params[Param]).replace(`*${Param}`, encodeURIComponent(Params[Param]));

    }

    const Response = await fetch(RouteURL, Options).catch((Err) => Log("Error", `Request Error: ${Err}`));

    if (!Response?.ok) {  return null; }

    const ResponseJSON = await Response.json().catch(() => null);

    if (!ResponseJSON) return null;

    return ResponseJSON;

}

// Logging

export function Log(Severity: "Info" | "Warning" | "Error", Message: string): void {

    const LoggingFunction = (Severity == "Info" ? console.log : Severity == "Warning" ? console.warn : console.error);

    LoggingFunction(`[${Severity.toUpperCase()}] ${Message}`);

}

// Auth

export function GetUserDirPath(): string {

    return `${GlobalStorage.User?.Username || "Home"}`;

}

// Paths

export function GetDisplayablePath(Path: string): string {

    return Path.replace("/Store", "");

}

// Cache

export function AddItemToCache(CacheType: keyof typeof GlobalStorage.Cache, Key: string, Item: any, Expiry: number = -1): void { // -1 means never expire as default

    GlobalStorage.Cache[CacheType][Key] = { Data: Item, Expiry };

}

export function WatchCache(): void {

    setInterval(() => {

        Log("Info", "Cache: Running clear operation...");

        for (const CacheType in GlobalStorage.Cache) {

            const Cache = GlobalStorage.Cache[CacheType as keyof typeof GlobalStorage.Cache];

            for (const Key in Cache) {

                const CacheItem = Cache[Key];

                if (CacheItem.Expiry < Date.now() && CacheItem.Expiry != -1) {

                    Log("Info", `Cache: Removing ${Key} from ${CacheType}...`);
                    delete Cache[Key];

                }

            }
            
        }

    }, 60_000);

}

// Internal Fetchers (Uses Cache)

export function FetchDir(UID: string): BackendDir | null {

    return GlobalStorage.Cache.Dirs[UID]?.Data || null;

}

export function FetchFile(UID: string): BackendFile | null {

    return GlobalStorage.Cache.Files[UID]?.Data || null;

}