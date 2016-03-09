package resource

type Quantity struct {
	// Amount is public, so you can manipulate it if the accessor
	// functions are not sufficient.
	//Amount *inf.Dec

	Amount *string
	// Change Format at will. See the comment for Canonicalize for
	// more details.
	Format
}
// Format lists the three possible formattings of a quantity.
type Format string
