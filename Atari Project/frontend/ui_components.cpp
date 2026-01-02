/**************************************/
/*                                    */
/*   UI Components - C++ Rendering    */
/*     Frutiger Aero + Y2K Edition    */
/*           Programmed by            */
/*            Sertaç Ataç             */
/*            02.01.2026              */
/*                                    */
/**************************************/

/**************************************************/
/*                                                */
/*    This file demonstrates C++ rendering        */
/*    components that can be called from Go       */
/*    via CGO. These functions provide low-level  */
/*    drawing utilities optimized for the         */
/*    retro gaming aesthetic.                     */
/*                                                */
/**************************************************/

#include <math.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**************************************************/
/*                                                */
/*       COLOR DEFINITIONS - NEON PALETTE         */
/*                                                */
/**************************************************/

typedef struct {
    uint8_t r, g, b, a;
} RGBAColor;

/*       Primary neon green palette              */
static const RGBAColor NEON_GREEN         = {57, 255, 20, 255};
static const RGBAColor NEON_GREEN_DARK    = {30, 180, 10, 255};
static const RGBAColor NEON_GREEN_GLOW    = {57, 255, 20, 100};
static const RGBAColor NEON_GREEN_BRIGHT  = {100, 255, 100, 255};

/*       Frutiger Aero backgrounds               */
static const RGBAColor AERO_BACKGROUND    = {15, 25, 35, 255};
static const RGBAColor AERO_DARK_PANEL    = {20, 35, 50, 230};
static const RGBAColor AERO_GLOSS         = {255, 255, 255, 40};
static const RGBAColor AERO_SHADOW        = {0, 0, 0, 150};

/**************************************************/
/*                                                */
/*              MATH UTILITIES                    */
/*                                                */
/**************************************************/

/*      Smooth interpolation for animations      */
float lerp_cpp(float a, float b, float t) {
    return a + (b - a) * t;
}

/*   Ease-out cubic for buttery smooth transitions   */
float ease_out_cubic(float t) {
    float invT = 1.0f - t;
    return 1.0f - (invT * invT * invT);
}

/*       Ease-in-out for bounce effects          */
float ease_in_out_quad(float t) {
    return t < 0.5f 
        ? 2.0f * t * t 
        : 1.0f - powf(-2.0f * t + 2.0f, 2) / 2.0f;
}

/*       Sinusoidal pulse for glow effects       */
float pulse_glow(float time, float frequency, float minVal, float maxVal) {
    float t = (sinf(time * frequency) + 1.0f) / 2.0f;
    return minVal + t * (maxVal - minVal);
}

/**************************************************/
/*                                                */
/*          GEOMETRIC CALCULATIONS                */
/*                                                */
/**************************************************/

/*   Calculate points for a blade shape          */
/*   (Xbox 360 style curved panel)               */
typedef struct {
    float x1, y1;  /* Top-left                   */
    float x2, y2;  /* Top-right                  */
    float x3, y3;  /* Bottom-right               */
    float x4, y4;  /* Bottom-left                */
    float curve;   /* Curve amount for right edge*/
} BladeGeometry;

BladeGeometry calculate_blade_geometry(
    float x, float y, 
    float width, float height,
    float curveOffset,
    int isActive
) {
    BladeGeometry geo;
    
    /*          Standard rectangle base          */
    geo.x1 = x;
    geo.y1 = y;
    geo.x2 = x + width;
    geo.y2 = y;
    geo.x3 = x + width;
    geo.y3 = y + height;
    geo.x4 = x;
    geo.y4 = y + height;
    
    /*   Apply curve to right edge for active blade   */
    if (isActive) {
        geo.curve = curveOffset;
        geo.x2 += curveOffset * 0.5f;
        geo.x3 += curveOffset * 0.5f;
    } else {
        geo.curve = 0.0f;
    }
    
    return geo;
}

/**************************************************/
/*                                                */
/*           GRADIENT CALCULATIONS                */
/*                                                */
/**************************************************/

typedef struct {
    RGBAColor topLeft;
    RGBAColor topRight;
    RGBAColor bottomLeft;
    RGBAColor bottomRight;
} QuadGradient;

/*   Generate Frutiger Aero style gradient for panels   */
QuadGradient generate_aero_gradient(RGBAColor baseColor, int hasGloss) {
    QuadGradient grad;
    
    /*       Base gradient - darker at bottom    */
    grad.topLeft = baseColor;
    grad.topRight = baseColor;
    
    /*              Darken bottom                */
    grad.bottomLeft = (RGBAColor){
        (uint8_t)(baseColor.r * 0.7f),
        (uint8_t)(baseColor.g * 0.7f),
        (uint8_t)(baseColor.b * 0.7f),
        baseColor.a
    };
    grad.bottomRight = grad.bottomLeft;
    
    /*        Add gloss highlight to top         */
    if (hasGloss) {
        grad.topLeft.r = (uint8_t)fminf(grad.topLeft.r + 30, 255);
        grad.topLeft.g = (uint8_t)fminf(grad.topLeft.g + 30, 255);
        grad.topLeft.b = (uint8_t)fminf(grad.topLeft.b + 30, 255);
        grad.topRight = grad.topLeft;
    }
    
    return grad;
}

/*          Generate neon glow gradient          */
QuadGradient generate_neon_glow(RGBAColor neonColor, uint8_t intensity) {
    QuadGradient grad;
    
    /*           Center is brightest             */
    RGBAColor glowColor = {
        neonColor.r, neonColor.g, neonColor.b, intensity
    };
    
    /*        Edges fade to transparent          */
    RGBAColor fadeColor = {
        neonColor.r, neonColor.g, neonColor.b, 0
    };
    
    grad.topLeft = fadeColor;
    grad.topRight = fadeColor;
    grad.bottomLeft = glowColor;
    grad.bottomRight = glowColor;
    
    return grad;
}

/**************************************************/
/*                                                */
/*              ANIMATION STATE                   */
/*                                                */
/**************************************************/

typedef struct {
    float bladeTransition;    /* 0.0 to 1.0 for blade slide  */
    float cardHoverScale;     /* Current hover scale         */
    float targetCardScale;    /* Target hover scale          */
    float glowPulseTime;      /* Time accumulator for glow   */
    float scanlineY;          /* Current scanline Y position */
    int selectedCard;         /* Currently selected card     */
    int targetBlade;          /* Target blade index          */
} AnimationState;

AnimationState create_animation_state() {
    AnimationState state = {0};
    state.bladeTransition = 1.0f;
    state.cardHoverScale = 1.0f;
    state.targetCardScale = 1.0f;
    return state;
}

void update_animation_state(AnimationState* state, float deltaTime) {
    /*        Update blade transition            */
    if (state->bladeTransition < 1.0f) {
        state->bladeTransition += deltaTime * 3.0f;
        if (state->bladeTransition > 1.0f) {
            state->bladeTransition = 1.0f;
        }
    }
    
    /*     Update card hover with smooth lerp    */
    state->cardHoverScale = lerp_cpp(
        state->cardHoverScale, 
        state->targetCardScale, 
        deltaTime * 8.0f
    );
    
    /*            Update glow pulse              */
    state->glowPulseTime += deltaTime * 2.0f;
    
    /*            Update scanline                */
    state->scanlineY += deltaTime * 200.0f;
    if (state->scanlineY > 720.0f) {
        state->scanlineY = 0.0f;
    }
}

/**************************************************/
/*                                                */
/*        Y2K BUBBLE PARTICLE SYSTEM              */
/*                                                */
/**************************************************/

#define MAX_BUBBLES 20

typedef struct {
    float x, y;
    float radius;
    float speed;
    uint8_t alpha;
    int active;
} Bubble;

typedef struct {
    Bubble bubbles[MAX_BUBBLES];
    int count;
} BubbleSystem;

BubbleSystem create_bubble_system() {
    BubbleSystem sys = {0};
    sys.count = MAX_BUBBLES;
    
    for (int i = 0; i < MAX_BUBBLES; i++) {
        sys.bubbles[i].active = 1;
        sys.bubbles[i].x = (float)(rand() % 1280);
        sys.bubbles[i].y = (float)(rand() % 720);
        sys.bubbles[i].radius = 5.0f + (rand() % 20);
        sys.bubbles[i].speed = 20.0f + (rand() % 40);
        sys.bubbles[i].alpha = 20 + (rand() % 40);
    }
    
    return sys;
}

void update_bubbles(BubbleSystem* sys, float deltaTime, int screenHeight) {
    for (int i = 0; i < sys->count; i++) {
        if (sys->bubbles[i].active) {
            /*          Float upward             */
            sys->bubbles[i].y -= sys->bubbles[i].speed * deltaTime;
            
            /*         Wrap around at top        */
            if (sys->bubbles[i].y < -sys->bubbles[i].radius) {
                sys->bubbles[i].y = (float)screenHeight + sys->bubbles[i].radius;
                sys->bubbles[i].x = (float)(rand() % 1280);
            }
        }
    }
}

/**************************************************/
/*                                                */
/*          CARD LAYOUT CALCULATIONS              */
/*                                                */
/**************************************************/

typedef struct {
    float x, y;
    float width, height;
    float scale;
    int isSelected;
} CardLayout;

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
) {
    CardLayout layout;
    
    layout.x = startX + (float)index * (cardWidth + spacing) - scrollOffset;
    layout.y = startY;
    layout.width = cardWidth;
    layout.height = cardHeight;
    layout.isSelected = (index == selectedIndex);
    
    if (layout.isSelected) {
        layout.scale = hoverScale;
        /*   Offset to keep card centered when scaled   */
        layout.x -= (cardWidth * (hoverScale - 1.0f)) / 2.0f;
        layout.y -= (cardHeight * (hoverScale - 1.0f)) / 2.0f;
    } else {
        layout.scale = 1.0f;
    }
    
    return layout;
}

/**************************************************/
/*                                                */
/*              COLOR BLENDING                    */
/*                                                */
/**************************************************/

RGBAColor blend_colors(RGBAColor c1, RGBAColor c2, float t) {
    return (RGBAColor){
        (uint8_t)(c1.r + (c2.r - c1.r) * t),
        (uint8_t)(c1.g + (c2.g - c1.g) * t),
        (uint8_t)(c1.b + (c2.b - c1.b) * t),
        (uint8_t)(c1.a + (c2.a - c1.a) * t)
    };
}

RGBAColor apply_alpha(RGBAColor color, float alphaMultiplier) {
    return (RGBAColor){
        color.r, color.g, color.b,
        (uint8_t)(color.a * alphaMultiplier)
    };
}

RGBAColor brighten_color(RGBAColor color, float amount) {
    return (RGBAColor){
        (uint8_t)fminf(color.r + amount * 255.0f, 255),
        (uint8_t)fminf(color.g + amount * 255.0f, 255),
        (uint8_t)fminf(color.b + amount * 255.0f, 255),
        color.a
    };
}

/**************************************************/
/*                                                */
/*         EXPORTED FUNCTIONS FOR CGO             */
/*                                                */
/**************************************************/

/*          Get eased transition value           */
float GetEasedTransition(float t) {
    return ease_out_cubic(t);
}

/*       Get pulse value for glow effects        */
float GetGlowPulse(float time) {
    return pulse_glow(time, 2.0f, 0.5f, 1.0f);
}

/*   Calculate blade X offset for animation      */
float CalculateBladeOffset(
    int currentBlade, 
    int targetBlade, 
    float bladeWidth, 
    float transitionProgress
) {
    float easeT = ease_out_cubic(transitionProgress);
    float currentOffset = (float)currentBlade * bladeWidth;
    float targetOffset = (float)targetBlade * bladeWidth;
    return currentOffset + (targetOffset - currentOffset) * easeT;
}

/*    Calculate card scale with bounce effect    */
float CalculateCardScale(float currentScale, float targetScale, float deltaTime) {
    return lerp_cpp(currentScale, targetScale, deltaTime * 8.0f);
}

#ifdef __cplusplus
}
#endif
