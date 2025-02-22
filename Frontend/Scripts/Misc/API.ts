import { GlobalStorage } from "../Main";

import type { BackendDir, BackendFile, RenderedResource } from "./Structs";

import { MakeRequest } from "./Utils";

import Routes from "../../Routes.json";

export async function GetDirContents(DirPath: string): Promise<RenderedResource<{ Parent: BackendDir | null, Files: BackendFile[], Dirs: BackendDir[] }>> {

    const Response = await MakeRequest(Routes.GetFileList, null, {}, {

        "Path": DirPath

    })

    if (!Response?.JSON?.Parent) return { HTML: "", JSON: { Parent: null, Files: [], Dirs: [] } };

    const Subfiles = Response.JSON.Files satisfies BackendFile[];
    const Subdirs = Response.JSON.Dirs satisfies BackendDir[];
    const Parent = Response.JSON.Parent satisfies BackendDir;

    return { HTML: Response.HTML, JSON: { Parent, Files: Subfiles, Dirs: Subdirs } };

}

export async function NewDirectory(Name: string, Path: string, Private: boolean, Authorized: string[] = []): Promise<boolean> {

    const Response = await MakeRequest(Routes.NewDir, { Name, Path, Private });

    GlobalStorage.Browser.GoTo(Path);

    return Response?.JSON?.Success ?? false;

}