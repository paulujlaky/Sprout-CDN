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

    GlobalStorage.Browser.Refresh();

    return Response?.JSON?.UID ?? false;

}

export async function NewFile(Path: string, Uploaded: File, Private: boolean = false): Promise<BackendFile | null> {

    const DataToSend = new FormData();

    DataToSend.append("File", Uploaded);
    DataToSend.append("Path", Path);
    DataToSend.append("Private", Private.toString());

    const Response = await MakeRequest(Routes.NewFile, DataToSend, {}, {}, true); // IsFormData must be set to true

    return Response?.JSON ? Response.JSON satisfies BackendFile : null;
    
}

export async function MoveFile(OldFilePath: string, PathOfNewDir: string): Promise<boolean> {

    const Response = await MakeRequest(Routes.MoveFile, { OldPath: OldFilePath, NewPath: PathOfNewDir });

    return Response?.JSON?.UID ?? false;

}

export async function UpdateFileAccess(FilePath: string, Private: boolean, Authorized: string[]): Promise<boolean> {

    const Response = await MakeRequest(Routes.UpdateFileAccess, { Path: FilePath, Private, Authorized });

    return Response?.JSON?.UID ?? false;

}

export async function RenameFile(OldFilePath: string, NewFileName: string): Promise<boolean> {

    const Response = await MakeRequest(Routes.RenameFile, { Path: OldFilePath, Name: NewFileName });

    return Response?.JSON?.UID ?? false;

}

export async function RenameDir(OldDirPath: string, NewDirName: string): Promise<boolean> {

    const Response = await MakeRequest(Routes.RenameDir, { Path: OldDirPath, Name: NewDirName });

    return Response?.JSON?.UID ?? false;

}

export async function MoveDir(OldDirPath: string, PathOfNewDir: string): Promise<boolean> {

    const Response = await MakeRequest(Routes.MoveDir, { OldPath: OldDirPath, NewPath: PathOfNewDir });

    return Response?.JSON?.UID ?? false;

}

export async function UpdateDirAccess(DirPath: string, Private: boolean, Authorized: string[]): Promise<boolean> {

    const Response = await MakeRequest(Routes.UpdateDirAccess, { Path: DirPath, Private, Authorized });

    return Response?.JSON?.UID ?? false;

}

export async function DeleteFile(FilePath: string): Promise<boolean> {

    const Response = await MakeRequest(Routes.DeleteFile, { Path: FilePath });

    return Response?.JSON?.Success ?? false;

}

export async function DeleteDir(DirPath: string, DirName: string): Promise<boolean> {

    const Response = await MakeRequest(Routes.DeleteDir, { Path: DirPath, Name: DirName });

    return Response?.JSON?.Success ?? false;

}

export async function GetFileQRCode(FilePath: string): Promise<Blob | null> {

    const Response = await fetch(Routes.GetFileQR.URL, {

        method: Routes.GetFileQR.Method,
        body: JSON.stringify({ Path: FilePath })

    }).catch((err) => console.error(err)); // Will return null if fetch fails

    // Will return image

    return Response?.ok ? await Response?.blob() : null;
    
}