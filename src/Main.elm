module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)



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
    , tutorialLink : String
    }


type alias Model =
    List Effect


init : flags -> ( Model, Cmd Msg )
init _ =
    ( [ Effect "stars"
            320
            200
            "Starfield 2D"
            "stars/stars.go"
            "Efeito de estrelas em scroll horizontal, com velocidade proporcional a posição do cursor do mouse sobre o canvas."
            "Issue_02_Introduction_To_Computer_Graphics.shtml"
      , Effect "3d"
            320
            200
            "Starfield 3D"
            "stars3D/stars3D.go"
            "Efeito de estrelas em movimentação 3D."
            ""
      , Effect "crossfade"
            200
            320
            "Crossfade"
            "crossfade/crossfade.go"
            "Efeito de transição entre duas imagens, pixel a pixel. Demonstra o uso de interpolação linear."
            "Issue_03_Timer_Related_Issues.shtml"
      , Effect "plasma"
            320
            200
            "Plasma"
            "plasma/plasma.go"
            "Efeito de plasma aplicado sobre imagem, pixel a pixel. Demonstra combinação de pixels manual."
            "Issue_04_Per_Pixel_Control.shtml"
      , Effect "filter"
            320
            200
            "Filter"
            "filters/filters.go"
            "Efeito de fogo usando filtros de imagem."
            "Issue_05_Filters.shtml"
      , Effect "cyber1"
            320
            200
            "Cyber1 / Lerp"
            "cyber1/cyber1.go"
            "Efeito de interpolação linear usando pixels da imagem. (sem link para tutorial pois esse foi da minha cabeça mesmo rsrs)"
            ""
      , Effect "bifilter"
            320
            200
            "Bi-Linear filter"
            "bifilter/bifilter.go"
            "Efeito de mapa de distorção com filtro bilinear (toque/clique na tela para alternar as versões!)."
            "Issue_06_Bitmap_Distortion.shtml"
      , Effect "bump"
            320
            200
            "Bump Map"
            "bump/bump.go"
            "Efeito de bump map (mapa de 'alturas'), com dois spots de luz se movendo aleatoriamente."
            "Issue_07_Bump_Mapping.shtml"
      , Effect "mandelbrot"
            320
            200
            "Zoom Fractal"
            "mandelbrot/mandelbrot.go"
            "Efeito de zoom fractal que desenha o conjunto Mandelbrot dando zoom in & out"
            "Issue_08_Fractal_Zooming.shtml"
      , Effect "textmap"
            320
            200
            "Textura Mapeada"
            "textmap/textmap.go"
            "Efeito de textura mapeada animada (du,dv)"
            "Issue_09_Static_Texture_Mapping.shtml"
      , Effect "rotozoom"
            320
            200
            "Roto Zoom"
            "rotozoom/rotozoom.go"
            "Efeito rotação e zoom"
            "Issue_10_Roto-Zooming.shtml"
      , Effect "particles"
            320
            200
            "Particles"
            "particles/particles.go"
            "Efeito Rastro"
            "Issue_11_Particle_Systems.shtml"
      , Effect "span"
            320
            200
            "Span"
            "span/span.go"
            "Efeito Span Based Rendering"
            "Issue_12_Span_Based_Rendering.shtml"
      , Effect "polygon"
            320
            200
            "Polygon Engines"
            "polygon/polygon.go"
            "Efeito Polygon Engines (rotating donut)"
            "Issue_13_Polygon_Engines.shtml"
      , Effect "plane"
            320
            200
            "Perspective Correct Texture Mapping"
            "plane/plane.go"
            "Efeito mapeamento de textura com perspectiva correta"
            "Issue_14_Perspective_Correct_Texture_Mapping.shtml"
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


singleEffect : Effect -> Html Msg
singleEffect e =
    div [ class "flex-shrink-0 p-sm-2 w-50" ]
        [ div [ class "card" ]
            [ div [ class "card-header d-flex justify-content-between align-items-center" ]
                [ text e.title
                , a
                    [ href ("https://github.com/moniquelive/demoscenetuts/blob/main/internal/" ++ e.code)
                    , target "_blank"
                    ]
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
                [ p [ class "card-text" ]
                    [ text (e.description ++ " ")
                    , if String.isEmpty e.tutorialLink then
                        text ""

                      else
                        a [ href ("https://www.flipcode.com/archives/The_Art_of_Demomaking-" ++ e.tutorialLink) ]
                            [ text "link para o tutorial" ]
                    ]
                ]
            ]
        ]


view : Model -> Html Msg
view model =
    div [ class "effects d-flex flex-wrap p-2" ] (List.map singleEffect model)
