module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Parser
import Html.Parser.Util



-- MAIN


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }



-- MODEL


type alias Effect =
    { effect : String
    , width : Int
    , height : Int
    , title : String
    , code : String
    , description : String
    }


type alias Model =
    List Effect


init : flags -> ( Model, Cmd Msg )
init _ =
    ( [ { effect = "stars"
        , width = 320
        , height = 200
        , title = "Starfield 2D"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/stars/stars.go"
        , description =
            """
      Efeito de estrelas em scroll horizontal, com velocidade proporcional a posição do cursor do mouse sobre o canvas.
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_02_Introduction_To_Computer_Graphics.shtml'>link para tutorial</a>
      """
        }
      , { effect = "3d"
        , width = 320
        , height = 200
        , title = "Starfield 3D"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/stars3D/stars3D.go"
        , description = "Efeito de estrelas em movimentação 3D."
        }
      , { effect = "crossfade"
        , width = 200
        , height = 320
        , title = "Crossfade"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/crossfade/crossfade.go"
        , description =
            """
      Efeito de transição entre duas imagens, pixel a pixel. Demonstra o uso de interpolação linear.
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_03_Timer_Related_Issues.shtml'>link para tutorial</a>
      """
        }
      , { effect = "plasma"
        , width = 320
        , height = 200
        , title = "Plasma"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/plasma/plasma.go"
        , description =
            """
      Efeito de plasma aplicado sobre imagem, pixel a pixel. Demonstra combinação de pixels manual.
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_04_Per_Pixel_Control.shtml'>link para tutorial</a>
      """
        }
      , { effect = "filter"
        , width = 320
        , height = 200
        , title = "Filter"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/filters/filters.go"
        , description =
            """
      Efeito de fogo usando filtros de imagem.
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_05_Filters.shtml'>link para tutorial</a>
      """
        }
      , { effect = "cyber1"
        , width = 320
        , height = 200
        , title = "Cyber1 / Lerp"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/cyber1/cyber1.go"
        , description = "Efeito de interpolação linear usando pixels da imagem. <i>Sem link para tutorial pois esse foi da minha cabeça mesmo rsrs</i>"
        }
      , { effect = "bifilter"
        , width = 320
        , height = 200
        , title = "Bi-Linear filter"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/bifilter/bifilter.go"
        , description =
            """
      Efeito de mapa de distorção com filtro bilinear (toque/clique na tela para alternar as versões!).
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_06_Bitmap_Distortion.shtml'>link para tutorial</a>
      """
        }
      , { effect = "bump"
        , width = 320
        , height = 200
        , title = "Bump Map"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/bump/bump.go"
        , description =
            """
      Efeito de bump map (mapa de 'alturas'), com dois spots de luz se movendo aleatoriamente.
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_07_Bump_Mapping.shtml'>link para tutorial</a>
      """
        }
      , { effect = "mandelbrot"
        , width = 320
        , height = 200
        , title = "Zoom Fractal"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/mandelbrot/mandelbrot.go"
        , description =
            """
      Efeito de zoom fractal que desenha o conjunto Mandelbrot dando zoom in & out
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_08_Fractal_Zooming.shtml'>link para tutorial</a>
      """
        }
      , { effect = "textmap"
        , width = 320
        , height = 200
        , title = "Textura Mapeada"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/textmap/textmap.go"
        , description =
            """
      Efeito de textura mapeada animada (du,dv)
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_09_Static_Texture_Mapping.shtml'>link para tutorial</a>
      """
        }
      , { effect = "rotozoom"
        , width = 320
        , height = 200
        , title = "Roto Zoom"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/rotozoom/rotozoom.go"
        , description =
            """
      Efeito rotação e zoom
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_10_Roto-Zooming.shtml'>link para tutorial</a>
      """
        }
      , { effect = "particles"
        , width = 320
        , height = 200
        , title = "Particles"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/particles/particles.go"
        , description =
            """
      Efeito Rastro
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_11_Particle_Systems.shtml'>link para tutorial</a>
      """
        }
      , { effect = "span"
        , width = 320
        , height = 200
        , title = "Span"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/span/span.go"
        , description =
            """
      Efeito Span Based Rendering
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_12_Span_Based_Rendering.shtml'>link para tutorial</a>
      """
        }
      , { effect = "polygon"
        , width = 320
        , height = 200
        , title = "Polygon Engines"
        , code = "https://github.com/moniquelive/demoscenetuts/blob/main/internal/polygon/polygon.go"
        , description =
            """
      Efeito Polygon Engines (rotating donut)
      <a href='https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_13_Polygon_Engines.shtml'>link para tutorial</a>
      """
        }
      ]
    , Cmd.none
    )



-- UPDATE


type Msg
    = Int


update : Msg -> Model -> ( Model, Cmd Msg )
update _ model =
    ( model, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none



-- VIEW


textHtml : String -> List (Html.Html msg)
textHtml t =
    case Html.Parser.run t of
        Ok nodes ->
            Html.Parser.Util.toVirtualDom nodes

        Err _ ->
            []


singleEffect : Effect -> Html Msg
singleEffect e =
    div [ class "flex-shrink-0 p-sm-2 w-50" ]
        [ div [ class "card" ]
            [ div [ class "card-header d-flex justify-content-between align-items-center" ]
                [ text e.title
                , a [ href e.code, target "_blank" ]
                    [ text "src" ]
                , button
                    [ type_ "button"
                    , class "btn btn-primary"
                    , attribute "data-bs-toggle" "modal"
                    , attribute "data-bs-target" "#iframeModal"
                    , attribute "data-effect" e.effect
                    , attribute "data-width" (String.fromInt e.width)
                    , attribute "data-height" (String.fromInt e.height)
                    ]
                    [ text "Exibir Efeito" ]
                ]
            , div [ class "card-body" ]
                [ p [ class "card-text" ] (textHtml e.description) ]
            ]
        ]


view : Model -> Html Msg
view model =
    div [ class "effects d-flex flex-wrap p-2" ] (List.map singleEffect model)
