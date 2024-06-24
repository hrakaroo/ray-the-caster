# Ray The Caster

Simple Ray Caster in SDL

## General Info

I'm not a huge gamer, but I can remember how I felt when I first saw DOOM. I was
blown away by the immersive feeling and have been hooked on 3D graphics ever
since. While I normally like to focus on Ray Tracers, I was curious why people
kept claiming that DOOM (and Wolfenstein3D) were not "true 3D environments". 
After watching a couple of videos on the making of Wolfenstein3D I decided
to give a go to creating my own basic Ray Caster.

I am not hyper optimizing for memory or speed. While I am curious how hard it is
to get a basic 3D environment working I am favoring readability over
speed. I don't mind some basic optimizations, but I'm not going to go too
crazy with it.

I am also not really interested in texture maps, or complex lighting. I really just
want to create a basic space I can explore and then I'll likely call it good.

## Build

I decided to write this in Go (since that is what I use most of the time at
work) with SDL for drawing.  In order to build this you first need to 

`brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config`

and then you should be able to just build and run it.

## Status

I'm mostly done. I still have a fish-eye going on, but I'm not as interested in fixing
that. I've done what I set out to do. 

## TODO

* (DONE) Get basic character (box) moving around the screen
* (DONE) Render walls
* (DONE) Add direction line to character
* (DONE) Motion in direction
* (DONE) Draw single casting ray
* (DONE) Draw Field of View casting rays
* (DONE) Render 3D
* (DONE) 3D shading using angles

## Notes & Useful Links

* https://github.com/veandco/go-sdl2-examples/tree/master
* https://github.com/veandco/go-sdl2
* https://lodev.org/cgtutor/raycasting.html
* https://www.youtube.com/watch?v=gYRrGTC7GtA&ab_channel=3DSage
* https://www.youtube.com/watch?v=NbSee-XM7WA&ab_channel=javidx9

## Images

Ray casting in 2D
![Ray casting in 2D](./img1.png)
