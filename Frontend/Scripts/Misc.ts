// Window / URL

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

export async function MakeRequest(Route: { URL: string, Method: string }, Body?: any, Headers: { [key: string]: string } = {}, IsFormData: boolean = false): Promise<BackendResponse | null> {

    const Options: RequestInit = {

        method: Route.Method,

        headers: IsFormData ? Headers : {

            "Content-Type": "application/json",
            ...Headers

        },

        body: IsFormData ? Body : JSON.stringify(Body)

    };

    const Response = await fetch(Route.URL, Options).catch(() => null);

    if (!Response?.ok) {  return null; }

    const ResponseJSON = await Response.json().catch(() => null);

    if (!ResponseJSON) return null;

    return ResponseJSON;

}