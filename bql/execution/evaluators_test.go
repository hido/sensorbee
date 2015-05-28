package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
	"time"
)

type evalTest struct {
	input    tuple.Value
	expected tuple.Value
}

func TestEvaluators(t *testing.T) {
	testCases := getTestCases()

	for _, testCase := range testCases {
		testCase := testCase
		ast := testCase.ast
		Convey(fmt.Sprintf("Given the AST Expression %v", ast), t, func() {

			Convey("Then an Evaluator can be computed", func() {
				eval, err := ExpressionToEvaluator(ast)
				So(err, ShouldBeNil)

				for i, tc := range testCase.inputs {
					input, expected := tc.input, tc.expected

					Convey(fmt.Sprintf("And when applied to input %v [%v]", input, i), func() {
						actual, err := eval.Eval(input)

						Convey(fmt.Sprintf("Then the result should be %v", expected), func() {
							i++
							if expected == nil {
								So(err, ShouldNotBeNil)
							} else {
								So(err, ShouldBeNil)
								So(actual, ShouldResemble, expected)
							}
						})
					})
				}

			})
		})
	}
}

func getTestCases() []struct {
	ast    interface{}
	inputs []evalTest
} {
	now := time.Now()

	// these are all type combinations that are so incompatible that
	// they cannot be compared with respect to less/greater and also
	// cannot be added etc.
	incomparables := []evalTest{
		// not a map:
		{tuple.Int(17), nil},
		// keys not present:
		{tuple.Map{"x": tuple.Int(17)}, nil},
		// only left present => error
		{tuple.Map{"a": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.Int(17)}, nil},
		{tuple.Map{"a": tuple.Float(3.14)}, nil},
		{tuple.Map{"a": tuple.String("日本語")}, nil},
		{tuple.Map{"a": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)}}, nil},
		// only right present => error
		{tuple.Map{"b": tuple.Bool(true)}, nil},
		{tuple.Map{"b": tuple.Int(17)}, nil},
		{tuple.Map{"b": tuple.Float(3.14)}, nil},
		{tuple.Map{"b": tuple.String("日本語")}, nil},
		{tuple.Map{"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// null vs. *
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.Int(3)}, nil},
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.Float(3.14)}, nil},
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.String("hoge")}, nil},
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Null{},
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// bool vs *
		{tuple.Map{"a": tuple.Bool(true),
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.Bool(true),
			"b": tuple.Int(3)}, nil},
		{tuple.Map{"a": tuple.Bool(true),
			"b": tuple.Float(3.14)}, nil},
		{tuple.Map{"a": tuple.Bool(true),
			"b": tuple.String("hoge")}, nil},
		{tuple.Map{"a": tuple.Bool(true),
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Bool(true),
			"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.Bool(true),
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Bool(true),
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// int vs. *
		{tuple.Map{"a": tuple.Int(3),
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.Int(3),
			"b": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.Int(3),
			"b": tuple.String("hoge")}, nil},
		{tuple.Map{"a": tuple.Int(3),
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Int(3),
			"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.Int(3),
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Int(3),
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// float vs *
		{tuple.Map{"a": tuple.Float(3.14),
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.Float(3.14),
			"b": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.Float(3.14),
			"b": tuple.String("hoge")}, nil},
		{tuple.Map{"a": tuple.Float(3.14),
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Float(3.14),
			"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.Float(3.14),
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Float(3.14),
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// string vs *
		{tuple.Map{"a": tuple.String("hoge"),
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.String("hoge"),
			"b": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.String("hoge"),
			"b": tuple.Int(3)}, nil},
		{tuple.Map{"a": tuple.String("hoge"),
			"b": tuple.Float(3.14)}, nil},
		{tuple.Map{"a": tuple.String("hoge"),
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.String("hoge"),
			"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.String("hoge"),
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.String("hoge"),
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// blob vs *
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.Int(3)}, nil},
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.Float(3.14)}, nil},
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.String("hoge")}, nil},
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Blob("hoge"),
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// timestamp vs *
		{tuple.Map{"a": tuple.Timestamp(now),
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.Timestamp(now),
			"b": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.Timestamp(now),
			"b": tuple.Int(3)}, nil},
		{tuple.Map{"a": tuple.Timestamp(now),
			"b": tuple.Float(3.14)}, nil},
		{tuple.Map{"a": tuple.Timestamp(now),
			"b": tuple.String("hoge")}, nil},
		{tuple.Map{"a": tuple.Timestamp(now),
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Timestamp(now),
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Timestamp(now),
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// array vs *
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.Int(3)}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.Float(3.14)}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.String("hoge")}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Array{tuple.Int(2)},
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
		// map vs *
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.Null{}}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.Bool(true)}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.Int(3)}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.Float(3.14)}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.String("hoge")}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.Blob("hoge")}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.Timestamp(now)}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.Array{tuple.Int(2)}}, nil},
		{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
			"b": tuple.Map{"b": tuple.Int(3)}}, nil},
	}

	// we should check that every AST expression maps to
	// an evaluator with the correct behavior
	testCases := []struct {
		ast    interface{}
		inputs []evalTest
	}{
		// Literals should always be independent of the input data
		{parser.NumericLiteral{23},
			[]evalTest{
				{tuple.Int(17), tuple.Int(23)},
				{tuple.String(""), tuple.Int(23)},
			},
		},
		{parser.FloatLiteral{3.14},
			[]evalTest{
				{tuple.Int(17), tuple.Float(3.14)},
				{tuple.String(""), tuple.Float(3.14)},
			},
		},
		{parser.BoolLiteral{true},
			[]evalTest{
				{tuple.Int(17), tuple.Bool(true)},
				{tuple.String(""), tuple.Bool(true)},
			},
		},
		{parser.BoolLiteral{false},
			[]evalTest{
				{tuple.Int(17), tuple.Bool(false)},
				{tuple.String(""), tuple.Bool(false)},
			},
		},
		// Access to columns/keys should return the same values
		{parser.ColumnName{"a"},
			[]evalTest{
				// not a map:
				{tuple.Int(17), nil},
				// key not present:
				{tuple.Map{"x": tuple.Int(17)}, nil},
				// key present
				{tuple.Map{"a": tuple.Null{}}, tuple.Null{}},
				{tuple.Map{"a": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(17)}, tuple.Int(17)},
				{tuple.Map{"a": tuple.Float(3.14)}, tuple.Float(3.14)},
				{tuple.Map{"a": tuple.String("日本語")}, tuple.String("日本語")},
				{tuple.Map{"a": tuple.Blob("hoge")}, tuple.Blob("hoge")},
				{tuple.Map{"a": tuple.Timestamp(now)}, tuple.Timestamp(now)},
				{tuple.Map{"a": tuple.Array{tuple.Int(2)}}, tuple.Array{tuple.Int(2)}},
				{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)}}, tuple.Map{"b": tuple.Int(3)}},
			},
		},
		/// Combined operations
		// Or
		{parser.BinaryOpAST{parser.Or, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			[]evalTest{
				// not a map:
				{tuple.Int(17), nil},
				// keys not present:
				{tuple.Map{"x": tuple.Int(17)}, nil},
				// only left key present and evaluates to true => right is not necessary
				{tuple.Map{"a": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(17)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.String("日本語")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Blob("hoge")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Timestamp(now)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Array{tuple.Int(2)}}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)}}, tuple.Bool(true)},
				// only left key present and evaluates to false => error
				{tuple.Map{"a": tuple.Null{}}, nil},
				{tuple.Map{"a": tuple.Bool(false)}, nil},
				{tuple.Map{"a": tuple.Int(0)}, nil},
				{tuple.Map{"a": tuple.Float(0.0)}, nil},
				{tuple.Map{"a": tuple.String("")}, nil},
				{tuple.Map{"a": tuple.Blob("")}, nil},
				{tuple.Map{"a": tuple.Timestamp{}}, nil},
				{tuple.Map{"a": tuple.Array{}}, nil},
				{tuple.Map{"a": tuple.Map{}}, nil},
				// left key evalues to false and right to true => true
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Int(17)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.String("日本語")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Blob("hoge")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Timestamp(now)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Array{tuple.Int(2)}}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Map{"b": tuple.Int(3)}}, tuple.Bool(true)},
				// left key evalues to false and right to false => false
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Bool(false)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Int(0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Float(0.0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.String("")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Blob("")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Timestamp{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Array{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(0),
					"b": tuple.Map{}}, tuple.Bool(false)},
			},
		},
		// And
		{parser.BinaryOpAST{parser.And, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			[]evalTest{
				// not a map:
				{tuple.Int(17), nil},
				// keys not present:
				{tuple.Map{"x": tuple.Int(17)}, nil},
				// only left key present and evaluates to false => right is not necessary
				{tuple.Map{"a": tuple.Null{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Bool(false)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(0.0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.String("")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Blob("")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Timestamp{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Array{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Map{}}, tuple.Bool(false)},
				// only left key present and evaluates to true => error
				{tuple.Map{"a": tuple.Bool(true)}, nil},
				{tuple.Map{"a": tuple.Int(17)}, nil},
				{tuple.Map{"a": tuple.Float(3.14)}, nil},
				{tuple.Map{"a": tuple.String("日本語")}, nil},
				{tuple.Map{"a": tuple.Blob("hoge")}, nil},
				{tuple.Map{"a": tuple.Timestamp(now)}, nil},
				{tuple.Map{"a": tuple.Array{tuple.Int(2)}}, nil},
				{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)}}, nil},
				// left key evalues to true and right to true => true
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Int(17)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.String("日本語")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Blob("hoge")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Timestamp(now)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Array{tuple.Int(2)}}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Map{"b": tuple.Int(3)}}, tuple.Bool(true)},
				// left key evalues to true and right to false => false
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Bool(false)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Int(0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Float(0.0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.String("")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Blob("")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Timestamp{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Array{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Map{}}, tuple.Bool(false)},
			},
		},
		/// Comparison Operations
		// Equal
		{parser.BinaryOpAST{parser.Equal, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			[]evalTest{
				// not a map:
				{tuple.Int(17), nil},
				// keys not present:
				{tuple.Map{"x": tuple.Int(17)}, nil},
				// only left present => error
				{tuple.Map{"a": tuple.Bool(true)}, nil},
				{tuple.Map{"a": tuple.Int(17)}, nil},
				{tuple.Map{"a": tuple.Float(3.14)}, nil},
				{tuple.Map{"a": tuple.String("日本語")}, nil},
				{tuple.Map{"a": tuple.Blob("hoge")}, nil},
				{tuple.Map{"a": tuple.Timestamp(now)}, nil},
				{tuple.Map{"a": tuple.Array{tuple.Int(2)}}, nil},
				{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)}}, nil},
				// only right present => error
				{tuple.Map{"b": tuple.Bool(true)}, nil},
				{tuple.Map{"b": tuple.Int(17)}, nil},
				{tuple.Map{"b": tuple.Float(3.14)}, nil},
				{tuple.Map{"b": tuple.String("日本語")}, nil},
				{tuple.Map{"b": tuple.Blob("hoge")}, nil},
				{tuple.Map{"b": tuple.Timestamp(now)}, nil},
				{tuple.Map{"b": tuple.Array{tuple.Int(2)}}, nil},
				{tuple.Map{"b": tuple.Map{"b": tuple.Int(3)}}, nil},
				// left and right present and equal => true
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(17),
					"b": tuple.Int(17)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.String("日本語"),
					"b": tuple.String("日本語")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Blob("hoge"),
					"b": tuple.Blob("hoge")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(now)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Array{tuple.Int(2)},
					"b": tuple.Array{tuple.Int(2)}}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
					"b": tuple.Map{"b": tuple.Int(3)}}, tuple.Bool(true)},
				// left and right present and not equal => false
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Bool(false)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Int(0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Float(0.0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.String("")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Blob("")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Timestamp{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Array{}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Map{}}, tuple.Bool(false)},
			},
		},
		// Less
		{parser.BinaryOpAST{parser.Less, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => true
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, tuple.Bool(true)},
				// left and right present and comparable and equal => false
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(true)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Int(3)}, tuple.Bool(false)},
				/* At the moment, Int(3) != Float(3.0):
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.0),
					"b": tuple.Int(3)}, tuple.Bool(false)},*/
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hoge")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(now)}, tuple.Bool(false)},
				// left and right present and comparable and left is greater => false
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(false)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Int(2)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(4),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(3)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.15),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.String("hogee"),
					"b": tuple.String("hoge")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Timestamp(time.Now()),
					"b": tuple.Timestamp(now)}, tuple.Bool(false)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// LessOrEqual
		{parser.BinaryOpAST{parser.LessOrEqual, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => true
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, tuple.Bool(true)},
				// left and right present and comparable and equal => true
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Int(3)}, tuple.Bool(true)},
				/* At the moment, Int(3) != Float(3.0):
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.0)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.0),
					"b": tuple.Int(3)}, tuple.Bool(true)},*/
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hoge")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(now)}, tuple.Bool(true)},
				// left and right present and comparable and left is greater => false
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(false)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Int(2)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(4),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(3)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.15),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.String("hogee"),
					"b": tuple.String("hoge")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Timestamp(time.Now()),
					"b": tuple.Timestamp(now)}, tuple.Bool(false)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Greater
		{parser.BinaryOpAST{parser.Greater, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => false
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, tuple.Bool(false)},
				// left and right present and comparable and equal => false
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(true)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Int(3)}, tuple.Bool(false)},
				/* At the moment, Int(3) != Float(3.0):
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.0)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.0),
					"b": tuple.Int(3)}, tuple.Bool(false)},*/
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hoge")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(now)}, tuple.Bool(false)},
				// left and right present and comparable and left is greater => true
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(false)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Int(2)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(4),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(3)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.15),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.String("hogee"),
					"b": tuple.String("hoge")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Timestamp(time.Now()),
					"b": tuple.Timestamp(now)}, tuple.Bool(true)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// GreaterOrEqual
		{parser.BinaryOpAST{parser.GreaterOrEqual, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => false
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, tuple.Bool(false)},
				// left and right present and comparable and equal => true
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(true)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Int(3)}, tuple.Bool(true)},
				/* At the moment, Int(3) != Float(3.0):
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.0)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.0),
					"b": tuple.Int(3)}, tuple.Bool(true)},*/
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hoge")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(now)}, tuple.Bool(true)},
				// left and right present and comparable and left is greater => true
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(false)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Int(2)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(4),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(3)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Float(3.15),
					"b": tuple.Float(3.14)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.String("hogee"),
					"b": tuple.String("hoge")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Timestamp(time.Now()),
					"b": tuple.Timestamp(now)}, tuple.Bool(true)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// NotEqual
		{parser.BinaryOpAST{parser.NotEqual, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			[]evalTest{
				// not a map:
				{tuple.Int(17), nil},
				// keys not present:
				{tuple.Map{"x": tuple.Int(17)}, nil},
				// only left present => error
				{tuple.Map{"a": tuple.Bool(true)}, nil},
				{tuple.Map{"a": tuple.Int(17)}, nil},
				{tuple.Map{"a": tuple.Float(3.14)}, nil},
				{tuple.Map{"a": tuple.String("日本語")}, nil},
				{tuple.Map{"a": tuple.Blob("hoge")}, nil},
				{tuple.Map{"a": tuple.Timestamp(now)}, nil},
				{tuple.Map{"a": tuple.Array{tuple.Int(2)}}, nil},
				{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)}}, nil},
				// only right present => error
				{tuple.Map{"b": tuple.Bool(true)}, nil},
				{tuple.Map{"b": tuple.Int(17)}, nil},
				{tuple.Map{"b": tuple.Float(3.14)}, nil},
				{tuple.Map{"b": tuple.String("日本語")}, nil},
				{tuple.Map{"b": tuple.Blob("hoge")}, nil},
				{tuple.Map{"b": tuple.Timestamp(now)}, nil},
				{tuple.Map{"b": tuple.Array{tuple.Int(2)}}, nil},
				{tuple.Map{"b": tuple.Map{"b": tuple.Int(3)}}, nil},
				// left and right present and equal => false
				{tuple.Map{"a": tuple.Bool(true),
					"b": tuple.Bool(true)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Int(17),
					"b": tuple.Int(17)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.14)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.String("日本語"),
					"b": tuple.String("日本語")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Blob("hoge"),
					"b": tuple.Blob("hoge")}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(now)}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Array{tuple.Int(2)},
					"b": tuple.Array{tuple.Int(2)}}, tuple.Bool(false)},
				{tuple.Map{"a": tuple.Map{"b": tuple.Int(3)},
					"b": tuple.Map{"b": tuple.Int(3)}}, tuple.Bool(false)},
				// left and right present and not equal => true
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Bool(false)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Int(0)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Float(0.0)}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.String("")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Blob("")}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Timestamp{}}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Array{}}, tuple.Bool(true)},
				{tuple.Map{"a": tuple.Int(1),
					"b": tuple.Map{}}, tuple.Bool(true)},
			},
		},
		/// Computational Operations
		// Plus
		{parser.BinaryOpAST{parser.Plus, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and can be added
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Int(5)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Float(float64(3) + 3.14)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Float(3.14 + float64(4))},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Float(3.14 + 3.15)},
				// left and right present and cannot be added
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, nil},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, nil},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Minus
		{parser.BinaryOpAST{parser.Minus, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and can be subtracted
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Int(-1)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Float(float64(3) - 3.14)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Float(3.14 - float64(4))},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Float(float64(3.14) - 3.15)},
				// left and right present and cannot be subtracted
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, nil},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, nil},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Multiply
		{parser.BinaryOpAST{parser.Multiply, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and can be multiplied
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Int(6)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Float(float64(3) * 3.14)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Float(3.14 * float64(4))},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Float(float64(3.14) * 3.15)},
				// left and right present and cannot be multiplied
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, nil},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, nil},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Divide
		{parser.BinaryOpAST{parser.Divide, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and can be divided
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Int(0)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Float(float64(3) / 3.14)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Float(3.14 / float64(4))},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Float(float64(3.14) / 3.15)},
				// division by zero
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(0)}, nil},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(0)}, tuple.Float(math.Inf(1))},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(0)}, tuple.Float(math.Inf(1))},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(0)}, tuple.Float(math.Inf(1))},
				// left and right present and cannot be divided
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, nil},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, nil},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Modulo
		{parser.BinaryOpAST{parser.Modulo, parser.ColumnName{"a"}, parser.ColumnName{"b"}},
			append([]evalTest{
				// left and right present and can be moduled
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(3)}, tuple.Int(2)},
				{tuple.Map{"a": tuple.Int(3),
					"b": tuple.Float(3.14)}, tuple.Float(3.0)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Int(4)}, tuple.Float(3.14)},
				{tuple.Map{"a": tuple.Float(3.14),
					"b": tuple.Float(3.15)}, tuple.Float(3.14)},
				// modulo by zero
				{tuple.Map{"a": tuple.Int(2),
					"b": tuple.Int(0)}, nil},
				// TODO add a way to check for IsNaN()
				/*
					{tuple.Map{"a": tuple.Int(3),
						"b": tuple.Float(0)}, tuple.Float(math.NaN())},
					{tuple.Map{"a": tuple.Float(3.14),
						"b": tuple.Int(0)}, tuple.Float(math.NaN())},
					{tuple.Map{"a": tuple.Float(3.14),
						"b": tuple.Float(0)}, tuple.Float(math.NaN())},
				*/
				// left and right present and cannot be moduled
				{tuple.Map{"a": tuple.Bool(false),
					"b": tuple.Bool(true)}, nil},
				{tuple.Map{"a": tuple.String("hoge"),
					"b": tuple.String("hogee")}, nil},
				{tuple.Map{"a": tuple.Timestamp(now),
					"b": tuple.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
	}
	return testCases
}
