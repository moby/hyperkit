package gist6545684_test

import (
	"github.com/shurcooL/go-goon"
	"github.com/shurcooL/go/gists/gist6545684"
)

func Example() {
	p, err := gist6545684.ReadGpcFile("./testdata/test_orientation.wwl")
	if err != nil {
		panic(err)
	}

	goon.Dump(p)

	// Output:
	// (gist6545684.Polygon)(gist6545684.Polygon{
	// 	Contours: ([]gist6545684.Contour)([]gist6545684.Contour{
	// 		(gist6545684.Contour)(gist6545684.Contour{
	// 			Vertices: ([]mgl64.Vec2)([]mgl64.Vec2{
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-210),
	// 					(float64)(-210),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(210),
	// 					(float64)(-210),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(210),
	// 					(float64)(210),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-210),
	// 					(float64)(210),
	// 				}),
	// 			}),
	// 		}),
	// 		(gist6545684.Contour)(gist6545684.Contour{
	// 			Vertices: ([]mgl64.Vec2)([]mgl64.Vec2{
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(180),
	// 					(float64)(180),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(180),
	// 					(float64)(120),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(120),
	// 					(float64)(180),
	// 				}),
	// 			}),
	// 		}),
	// 		(gist6545684.Contour)(gist6545684.Contour{
	// 			Vertices: ([]mgl64.Vec2)([]mgl64.Vec2{
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(0),
	// 					(float64)(180),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-15),
	// 					(float64)(150),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(15),
	// 					(float64)(150),
	// 				}),
	// 			}),
	// 		}),
	// 		(gist6545684.Contour)(gist6545684.Contour{
	// 			Vertices: ([]mgl64.Vec2)([]mgl64.Vec2{
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-150),
	// 					(float64)(150),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-134),
	// 					(float64)(150),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-134),
	// 					(float64)(166),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-150),
	// 					(float64)(166),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-150),
	// 					(float64)(150),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-166),
	// 					(float64)(150),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-166),
	// 					(float64)(134),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-150),
	// 					(float64)(134),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-150),
	// 					(float64)(118),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-134),
	// 					(float64)(118),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-134),
	// 					(float64)(134),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-150),
	// 					(float64)(134),
	// 				}),
	// 				(mgl64.Vec2)(mgl64.Vec2{
	// 					(float64)(-150),
	// 					(float64)(150),
	// 				}),
	// 			}),
	// 		}),
	// 	}),
	// })
}
