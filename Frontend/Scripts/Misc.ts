// Window / URL

export function GetURLParameter(Name: string): string | null {

    const URLInstance = new URL(window.location.href);
    return URLInstance.searchParams.get(Name);

}

// Network

interface BackendResponse {

    Message: string;

    JSON: any;
    HTML: string;

}

export async function MakeRequest(Type: "GET" | "POST" | "DELETE" | "PATCH", URL: string, Content: { Body: any, IsFormData?: boolean }, Headers: { [key: string]: string }): Promise<BackendResponse | null> {

    const Options: RequestInit = {

        method: Type,

        headers: Content.IsFormData ? Headers : {

            "Content-Type": "application/json",
            ...Headers

        },

        body: Content.IsFormData ? (Content.IsFormData ? Content.Body : JSON.stringify(Content.Body)) : null // FormData is not stringifiable

    };

    const Response = await fetch(URL, Options).catch(() => null);

    if (!Response?.ok) {  return null; }

    const ResponseJSON = await Response.json().catch(() => null);

    if (!ResponseJSON) return null;

    return {

        Message: ResponseJSON.Message,
        JSON: ResponseJSON,

        HTML: ""

    };

}