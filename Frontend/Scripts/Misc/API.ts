import { GlobalStorage } from "../Main";

import type { BackendDir, BackendFile, RenderedResource } from "./Structs";

import { AddItemToCache, MakeRequest } from "./Utils";

import Routes from "../../Routes.json";

// Misc

export async function GetDirContents(DirPath: string): Promise<RenderedResource<{ Parent: BackendDir | null, Files: BackendFile[], Dirs: BackendDir[] }>> {

    const Response = await MakeRequest(Routes.GetFileList, null, {}, {

        "Path": DirPath

    })

    if (!Response?.JSON?.Parent) return { HTML: "", JSON: { Parent: null, Files: [], Dirs: [] } };

    const Subfiles = Response.JSON.Files satisfies BackendFile[];
    const Subdirs = Response.JSON.Dirs satisfies BackendDir[];
    const Parent = Response.JSON.Parent satisfies BackendDir;

    for (const File of Subfiles) {

        AddItemToCache("Files", File.UID, File);

    }

    for (const Dir of Subdirs) {

        AddItemToCache("Dirs", Dir.UID, Dir);

    }

    return { HTML: Response.HTML, JSON: { Parent, Files: Subfiles, Dirs: Subdirs } };

}

// Creation Ops

export async function NewDirectory(Name: string, Path: string, Private: boolean, Authorized: string[] = []): Promise<boolean> {

    const Response = await MakeRequest(Routes.NewDir, { Name, Path, Private });

    GlobalStorage.Browser.GoTo(Path);

    return Response?.JSON?.UID ?? false;

}

export async function NewFile(Path: string, Uploaded: File, Private: boolean = false): Promise<boolean> {

    const DataToSend = new FormData();

    DataToSend.append("File", Uploaded);
    DataToSend.append("Path", Path);
    DataToSend.append("Private", Private.toString());

    const Response = await MakeRequest(Routes.NewFile, DataToSend, {}, {}, true); // IsFormData must be set to true

    return Response?.JSON?.UID ?? false; // Existence of uploaded file UID means success
    
}