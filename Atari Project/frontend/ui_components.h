/**************************************/
/*                                    */
/*       UI Components Header         */
/*     Frutiger Aero + Y2K Edition    */
/*           Programmed by            */
/*            Sertaç Ataç             */
/*            02.01.2026              */
/*                                    */
/**************************************/

#ifndef UI_COMPONENTS_H
#define UI_COMPONENTS_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**************************************************/
/*                                                */
/*            TYPE DEFINITIONS                    */
/*                                                */
/**************************************************/

/*              RGBA Color structure              */
typedef struct {
    uint8_t r, g, b, a;
} RGBAColor;

/*             Blade panel geometry               */
typedef struct {
    float x1, y1;  /* Top-left                   */
    float x2, y2;  /* Top-right                  */
    float x3, y3;  /* Bottom-right               */
    float x4, y4;  /* Bottom-left                */
    float curve;   /* Curve for right edge       */
} BladeGeometry;

/*         Quad gradient for Aero effects         */
typedef struct {
    RGBAColor topLeft;
    RGBAColor topRight;
    RGBAColor bottomLeft;
    RGBAColor bottomRight;
} QuadGradient;

/*           Animation state container            */
typedef struct {
    float bladeTransition;    /* 0.0 to 1.0 for blade slide  */
    float cardHoverScale;     /* Current hover scale         */
    float targetCardScale;    /* Target hover scale          */
    float glowPulseTime;      /* Time accumulator for glow   */
    float scanlineY;          /* Current scanline Y position */
    int selectedCard;         /* Currently selected card     */
    int targetBlade;          /* Target blade index          */
} AnimationState;

/*            Card layout information             */
typedef struct {
    float x, y;
    float width, height;
    float scale;
    int isSelected;
} CardLayout;

/**************************************************/
/*                                                */
/*              COLOR CONSTANTS                   */
/*                                                */
/**************************************************/

/*         Primary neon green palette             */
static const RGBAColor NEON_GREEN         = {57, 255, 20, 255};
static const RGBAColor NEON_GREEN_DARK    = {30, 180, 10, 255};
static const RGBAColor NEON_GREEN_GLOW    = {57, 255, 20, 100};
static const RGBAColor NEON_GREEN_BRIGHT  = {100, 255, 100, 255};

/*         Frutiger Aero backgrounds              */
static const RGBAColor AERO_BACKGROUND    = {15, 25, 35, 255};
static const RGBAColor AERO_DARK_PANEL    = {20, 35, 50, 230};
static const RGBAColor AERO_GLOSS         = {255, 255, 255, 40};
static const RGBAColor AERO_SHADOW        = {0, 0, 0, 150};

/**************************************************/
/*                                                */
/*              MATH UTILITIES                    */
/*                                                */
/**************************************************/

float lerp_cpp(float a, float b, float t);
float ease_out_cubic(float t);
float ease_in_out_quad(float t);
float pulse_glow(float time, float frequency, float minVal, float maxVal);

/**************************************************/
/*                                                */
/*           GEOMETRY FUNCTIONS                   */
/*                                                */
/**************************************************/

BladeGeometry calculate_blade_geometry(
    float x, float y, 
    float width, float height,
    float curveOffset,
    int isActive
);

/**************************************************/
/*                                                */
/*           GRADIENT FUNCTIONS                   */
/*                                                */
/**************************************************/

QuadGradient generate_aero_gradient(RGBAColor baseColor, int hasGloss);
QuadGradient generate_neon_glow(RGBAColor neonColor, uint8_t intensity);

/**************************************************/
/*                                                */
/*          ANIMATION FUNCTIONS                   */
/*                                                */
/**************************************************/

AnimationState create_animation_state(void);
void update_animation_state(AnimationState* state, float deltaTime);

/**************************************************/
/*                                                */
/*          CARD LAYOUT FUNCTIONS                 */
/*                                                */
/**************************************************/

CardLayout calculate_card_layout(
    int index,
    int selectedIndex,
    float hoverScale,
    float scrollOffset,
    float startX,
    float startY,
    float cardWidth,
    float cardHeight,
    float spacing
);

/**************************************************/
/*                                                */
/*         COLOR BLENDING FUNCTIONS               */
/*                                                */
/**************************************************/

RGBAColor blend_colors(RGBAColor c1, RGBAColor c2, float t);
RGBAColor apply_alpha(RGBAColor color, float alphaMultiplier);
RGBAColor brighten_color(RGBAColor color, float amount);

/**************************************************/
/*                                                */
/*         CGO EXPORTED FUNCTIONS                 */
/*                                                */
/**************************************************/

float GetEasedTransition(float t);
float GetGlowPulse(float time);
float CalculateBladeOffset(
    int currentBlade, 
    int targetBlade, 
    float bladeWidth, 
    float transitionProgress
);
float CalculateCardScale(float currentScale, float targetScale, float deltaTime);

#ifdef __cplusplus
}
#endif

#endif /* UI_COMPONENTS_H */
