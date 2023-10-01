package main

import (
	"fmt"
	"os"
	"time"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	// "github.com/damianpeaf/OLC2_Proyecto2_202110568/cst"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/errors"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/visitor"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/repl"

	"github.com/antlr4-go/antlr/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/compile", func(c *fiber.Ctx) error {

		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Println("Recovered in f", r)
		// 	}
		// }()

		startTime := time.Now()

		// resultChannel := make(chan string)

		// go func() {
		// 	resultChannel <- cst.CstReport(c.FormValue("code"))
		// }()

		code := c.FormValue("code")

		lexicalErrorListener := errors.NewLexicalErrorListener()
		lexer := compiler.NewTSwiftLexer(antlr.NewInputStream(code))

		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(lexicalErrorListener)

		stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

		parser := compiler.NewTSwiftLanguage(stream)
		parser.BuildParseTrees = true

		syntaxErrorListener := errors.NewSyntaxErrorListener(lexicalErrorListener.ErrorTable)
		parser.RemoveErrorListeners()
		parser.SetErrorHandler(errors.NewCustomErrorStrategy())
		parser.AddErrorListener(syntaxErrorListener)

		tree := parser.Program()

		// ♻ Recycle 'proyecto1' repl
		// TODO: sema checker should be a visitor instead of using the repl
		dclVisitor := repl.NewDclVisitor(syntaxErrorListener.ErrorTable)
		dclVisitor.Visit(tree)
		replVisitor := repl.NewVisitor(dclVisitor)
		replVisitor.Visit(tree)
		intepretationEndTime := time.Now()
		// ♻ End Recycle 'proyecto1' repl

		// begin ir generation
		irVisitor := visitor.NewIrVisitor()
		irVisitor.Visit(tree)
		fmt.Println("IR generation finished")
		irResult := irVisitor.Factory.String()
		fmt.Println(irResult)

		// cstReport := <-resultChannel

		reportEndTime := time.Now()
		fmt.Println("Interpretation finished")

		fmt.Println("Interpretation time:", intepretationEndTime.Sub(startTime))
		fmt.Println("Total (with report) time:", reportEndTime.Sub(intepretationEndTime))
		fmt.Println("")

		return c.JSON(struct {
			Errors     []repl.Error     `json:"errors"`
			Output     string           `json:"output"`
			CSTSvg     string           `json:"cstSvg"`
			ScopeTrace repl.ReportTable `json:"scopeTrace"`
			C3D        string           `json:"c3d"`
		}{
			Errors: replVisitor.ErrorTable.Errors,
			Output: replVisitor.Console.GetOutput(),
			// CSTSvg:     cstReport,
			ScopeTrace: replVisitor.ScopeTrace.Report(),
			C3D:        irResult,
		})

	})

	app.Listen(getPort())

}
