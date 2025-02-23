import $ from 'jquery';
import { GlobalStorage } from '../Main';

const RelevantElements = {

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

        alert("New Folder Button Clicked");

    });

    // Only inlineFiles and InlineDirs

    $(document).on("click", ".InlineFile, .InlineDir", (Event) => {

        const Target = $(Event.target);

        if (Target.attr("URL")) {

            window.open(Target.attr("URL") || "", "_blank");

        }
        

        

    });
    

}