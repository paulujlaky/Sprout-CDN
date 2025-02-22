import $ from "jquery";

import type { FullDirectory } from "../Models.ts/Dir";

import { WriteURL } from "../Misc/Utils";

const RelevantElements = {

    FileList: $(".DashMainFileList"),

}

export function RenderDirectory(Dir: FullDirectory): void {

    RelevantElements.FileList.html(Dir.Contents.HTML);

    WriteURL(Dir.Data.Path.replace("Store", "Dash"));

}