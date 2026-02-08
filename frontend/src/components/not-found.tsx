import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Link, useRouter } from "@tanstack/react-router";
import { ArrowLeft, Compass, Home, MapPinOff, Search } from "lucide-react";
import { useMemo } from "react";

export function NotFound() {
  const router = useRouter();

  // Generate floating particles for visual interest
  const particles = useMemo(() => {
    return Array.from({ length: 12 }).map((_, i) => ({
      id: i,
      size: Math.random() * 4 + 2,
      x: Math.random() * 100,
      y: Math.random() * 100,
      duration: Math.random() * 10 + 15,
      delay: Math.random() * 5,
    }));
  }, []);

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 relative overflow-hidden">
      {/* Floating particles */}
      <div className="absolute inset-0 pointer-events-none overflow-hidden">
        {particles.map((particle) => (
          <div
            key={particle.id}
            className="absolute rounded-full bg-primary/20"
            style={{
              width: particle.size,
              height: particle.size,
              left: `${particle.x}%`,
              top: `${particle.y}%`,
              animation: `float ${particle.duration}s ease-in-out ${particle.delay}s infinite`,
            }}
          />
        ))}
      </div>

      {/* Radial gradient overlay */}
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,transparent_0%,hsl(var(--background))_70%)] pointer-events-none" />

      <div className="w-full max-w-lg relative z-10">
        {/* Large 404 Display */}
        <div className="text-center mb-8 relative">
          <div className="relative inline-block">
            {/* Glowing effect behind text */}
            <div className="absolute inset-0 blur-3xl opacity-10 bg-primary" style={{ transform: "scale(1.5)" }} />

            {/* Main 404 text */}
            <div className="relative flex items-center justify-center gap-2">
              <span className="text-[10rem] font-mono font-black leading-none tracking-tighter text-foreground/5 select-none">
                4
              </span>
              <div className="relative">
                <Compass className="w-28 h-28 text-primary/20 animate-[spin_20s_linear_infinite]" />
                <MapPinOff className="w-12 h-12 text-primary absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
              </div>
              <span className="text-[10rem] font-mono font-black leading-none tracking-tighter text-foreground/5 select-none">
                4
              </span>
            </div>
          </div>
        </div>

        {/* Main Card */}
        <Card className="relative overflow-hidden border-2 border-border/50 backdrop-blur-sm">
          {/* Decorative gradient */}
          <div className="absolute inset-0 bg-gradient-to-br from-muted/30 to-transparent pointer-events-none" />

          {/* Corner accents */}
          <div className="absolute top-0 left-0 w-16 h-16 border-l-2 border-t-2 border-primary/20 rounded-tl-lg pointer-events-none" />
          <div className="absolute bottom-0 right-0 w-16 h-16 border-r-2 border-b-2 border-primary/20 rounded-br-lg pointer-events-none" />

          <CardContent className="relative p-8">
            {/* Icon */}
            <div className="flex justify-center mb-6">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-muted/50 to-muted/20 border border-border/50 flex items-center justify-center">
                <Search className="w-8 h-8 text-muted-foreground" />
              </div>
            </div>

            {/* Title & Description */}
            <div className="text-center mb-8">
              <h1 className="text-2xl font-bold tracking-tight mb-2 font-mono">
                Page Not Found
              </h1>
              <p className="text-muted-foreground text-sm leading-relaxed max-w-sm mx-auto">
                The page you're looking for doesn't exist or has been moved to a different location.
              </p>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-col sm:flex-row gap-3 justify-center">
              <Button
                variant="outline"
                onClick={() => router.history.back()}
                className="gap-2"
              >
                <ArrowLeft className="w-4 h-4" />
                Go Back
              </Button>

              <Button variant="glow" asChild className="gap-2">
                <Link to="/dashboard">
                  <Home className="w-4 h-4" />
                  Dashboard
                </Link>
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Helpful hint */}
        <p className="text-center text-xs text-muted-foreground mt-6">
          Lost? Try using the navigation menu or search for what you need.
        </p>
      </div>
    </div>
  );
}
