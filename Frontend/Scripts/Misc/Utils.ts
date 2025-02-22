// Window / URL

import { GlobalStorage } from "../Main";

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

    window.history.pushState({}, "", URL);

}

// Network

interface BackendResponse {

    Message: string;

    JSON: any;
    HTML: string;

}

export async function MakeRequest(Route: { URL: string, Method: string }, Body?: any, Headers: { [key: string]: string } = {}, Params: { [key: string]: string } = {}, IsFormData: boolean = false): Promise<BackendResponse | null> {

    const Options: RequestInit = {

        method: Route.Method,

        headers: IsFormData ? Headers : {

            "Content-Type": "application/json",
            ...Headers

        },

        body: IsFormData ? Body : JSON.stringify(Body)

    };

    for (const Param in Params) {

        Route.URL = Route.URL.replace(`:${Param}`, Params[Param]).replace(`*${Param}`, encodeURIComponent(Params[Param]));

    }

    const Response = await fetch(Route.URL, Options).catch(() => null);

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

    return `/${GlobalStorage.User?.Username || "Home"}`;

}