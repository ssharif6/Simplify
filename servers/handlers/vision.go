package handlers
import (
	"fmt"
	"io"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/net/context"
)

func init() {
	// Refer to these functions so that goimports is happy before boilerplate is inserted.
	_ = context.Background()
	_ = vision.ImageAnnotatorClient{}
	_ = os.Open
}

// detectLabels gets labels from the Vision API for an image at the given file path.
func DetectLabels(file string) ([]string, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}
	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}

	stringSlice := make([]string, 0)
	//
	//for _, c := range texts {
	//	stringSlice = append(stringSlice, c.Description)
	//}
	if len(annotations) == 0 {
	} else {

		for _, annotation := range annotations {
			stringSlice = append(stringSlice, annotation.Description)
		}
	}

	return stringSlice, nil
}

// detectLandmarks gets landmarks from the Vision API for an image at the given file path.
func detectLandmarks(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectLandmarks(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No landmarks found.")
	} else {
		fmt.Fprintln(w, "Landmarks:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// detectText gets text from the Vision API for an image at the given file path.
func DetectText(file string) ([]string, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}

	stringSlice := make([]string, 0)
	//
	//for _, c := range texts {
	//	stringSlice = append(stringSlice, c.Description)
	//}
	if len(annotations) == 0 {
	} else {

		for _, annotation := range annotations {
				stringSlice = append(stringSlice, annotation.Description)
		}
	}

	//r := &ImageResponseObject{
	//	Texts: stringSlice,
	//}

	return stringSlice, nil
}

// detectDocumentText gets the full document text from the Vision API for an image at the given file path.
func detectDocumentText(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return err
	}

	if annotation == nil {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Text:")
		fmt.Fprintf(w, "%q\n", annotation.Text)
	}

	return nil
}

// detectProperties gets image properties from the Vision API for an image at the given file path.
func detectProperties(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	props, err := client.DetectImageProperties(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Dominant colors:")
	for _, quantized := range props.DominantColors.Colors {
		color := quantized.Color
		r := int(color.Red) & 0xff
		g := int(color.Green) & 0xff
		b := int(color.Blue) & 0xff
		fmt.Fprintf(w, "%2.1f%% - #%02x%02x%02x\n", quantized.PixelFraction*100, r, g, b)
	}

	return nil
}

// detectCropHints gets suggested croppings the Vision API for an image at the given file path.
func detectCropHints(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	res, err := client.CropHints(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Crop hints:")
	for _, hint := range res.CropHints {
		for _, v := range hint.BoundingPoly.Vertices {
			fmt.Fprintf(w, "(%d,%d)\n", v.X, v.Y)
		}
	}

	return nil
}

// detectSafeSearch gets image properties from the Vision API for an image at the given file path.
func detectSafeSearch(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	props, err := client.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Safe Search properties:")
	fmt.Fprintln(w, "Adult:", props.Adult)
	fmt.Fprintln(w, "Medical:", props.Medical)
	fmt.Fprintln(w, "Spoofed:", props.Spoof)
	fmt.Fprintln(w, "Violence:", props.Violence)

	return nil
}

// detectWeb gets image properties from the Vision API for an image at the given file path.
func detectWeb(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	web, err := client.DetectWeb(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Web properties:")
	if len(web.FullMatchingImages) != 0 {
		fmt.Fprintln(w, "\tFull image matches:")
		for _, full := range web.FullMatchingImages {
			fmt.Fprintf(w, "\t\t%s\n", full.Url)
		}
	}
	if len(web.PagesWithMatchingImages) != 0 {
		fmt.Fprintln(w, "\tPages with this image:")
		for _, page := range web.PagesWithMatchingImages {
			fmt.Fprintf(w, "\t\t%s\n", page.Url)
		}
	}
	if len(web.WebEntities) != 0 {
		fmt.Fprintln(w, "\tEntities:")
		for _, entity := range web.WebEntities {
			fmt.Fprintf(w, "\t\t%-12s %s\n", entity.EntityId, entity.Description)
		}
	}

	return nil
}

// detectLogos gets logos from the Vision API for an image at the given file path.
func detectLogos(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectLogos(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No logos found.")
	} else {
		fmt.Fprintln(w, "Logos:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

