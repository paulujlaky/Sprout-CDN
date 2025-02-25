import $ from 'jquery';
import { GlobalStorage } from '../Main';
import { FetchFile } from '../Misc/Utils';

const RelevantElements = {

    Document: $(document),

    Buttons: {

        NewFile: $(".DashFooterButton.NewFile"),
        NewFolder: $(".DashFooterButton.NewFolder"),

    }    

}

export function WatchPageInteractions(): void {

    RelevantElements.Buttons.NewFile.on("click", () => {

        // Prompt user to upload new file

        const FileInput = document.createElement("input");

        FileInput.type = "file";
        FileInput.accept = "*/*";
        
        FileInput.click();

        FileInput.onchange = async () => {

            const File = (FileInput.files || [])[0];

            if (!File) { return; }

            // Upload file

            await GlobalStorage.Browser.Current?.Upload(File);

        };
        
    });

    RelevantElements.Buttons.NewFolder.on("click", () => {

        
        

    });

    const InlineFileInteractionRouter = (Event: JQuery.MouseEventBase | JQuery.TouchEventBase): void => {

        const Target = $(Event.target).closest(".InlineFile");

        if ($(Event.target).closest(".InlineFileActions ").length > 0) { return; }

        const UID = Target.attr("UID");

        if (!UID) { return; }

        const RelevantFile = FetchFile(UID);

        console.log(RelevantFile);

    };
    
    RelevantElements.Document.on("click", InlineFileInteractionRouter);
        
}