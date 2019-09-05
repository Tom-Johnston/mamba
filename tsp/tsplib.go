package tsp

import (
	"fmt"
	"io"
	"text/tabwriter"
)

//LIB writes to w a the TSP problem with the given integer weights in a TSPLIB compatible format.
//The number of vertices (also called the dimension) is given by n and the function weights returns the weight of the edge from i to j (which is the same as j to i).
func LIB(w io.Writer, n int, weights func(i, j int) int) (err error) {
	_, err = io.WriteString(w, "TYPE: TSP\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "DIMENSION: %d\n", n)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, "DISPLAY_DATA_TYPE: NO_DISPLAY\nEDGE_WEIGHT_TYPE: EXPLICIT\nEDGE_WEIGHT_FORMAT: LOWER_DIAG_ROW\nEDGE_WEIGHT_SECTION\n")
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 1, 1, ' ', tabwriter.AlignRight)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			fmt.Fprintf(tw, "%d\t", weights(i, j))
		}
		fmt.Fprint(tw, "0\t")
		fmt.Fprint(tw, "\n")
	}
	tw.Flush()
	_, err = io.WriteString(w, "EOF\n")
	if err != nil {
		return err
	}
	return nil
}
