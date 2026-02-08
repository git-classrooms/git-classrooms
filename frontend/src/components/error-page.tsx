import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Link, useRouter } from "@tanstack/react-router";
import {
  AlertTriangle,
  ArrowLeft,
  Ban,
  Home,
  Lock,
  RefreshCw,
  ServerCrash,
  ShieldX,
  Wifi,
  WifiOff,
} from "lucide-react";
import { useMemo } from "react";

interface ErrorConfig {
  icon: React.ElementType;
  title: string;
  description: string;
  iconColor: string;
  iconBg: string;
  accentColor: string;
}

const errorConfigs: Record<number, ErrorConfig> = {
  400: {
    icon: AlertTriangle,
    title: "Bad Request",
    description: "The server couldn't understand your request. Please check your input and try again.",
    iconColor: "text-warning",
    iconBg: "from-warning/20 to-warning/5",
    accentColor: "warning",
  },
  401: {
    icon: Lock,
    title: "Unauthorized",
    description: "You need to sign in to access this resource. Please authenticate and try again.",
    iconColor: "text-warning",
    iconBg: "from-warning/20 to-warning/5",
    accentColor: "warning",
  },
  403: {
    icon: ShieldX,
    title: "Access Denied",
    description: "You don't have permission to access this resource. Contact your administrator if you believe this is an error.",
    iconColor: "text-destructive",
    iconBg: "from-destructive/20 to-destructive/5",
    accentColor: "destructive",
  },
  404: {
    icon: Ban,
    title: "Not Found",
    description: "The page you're looking for doesn't exist or has been moved.",
    iconColor: "text-muted-foreground",
    iconBg: "from-muted/50 to-muted/20",
    accentColor: "muted",
  },
  408: {
    icon: WifiOff,
    title: "Request Timeout",
    description: "The server took too long to respond. Please check your connection and try again.",
    iconColor: "text-warning",
    iconBg: "from-warning/20 to-warning/5",
    accentColor: "warning",
  },
  500: {
    icon: ServerCrash,
    title: "Server Error",
    description: "Something went wrong on our end. Our team has been notified and is working on it.",
    iconColor: "text-destructive",
    iconBg: "from-destructive/20 to-destructive/5",
    accentColor: "destructive",
  },
  502: {
    icon: Wifi,
    title: "Bad Gateway",
    description: "We're having trouble connecting to our servers. Please try again in a moment.",
    iconColor: "text-destructive",
    iconBg: "from-destructive/20 to-destructive/5",
    accentColor: "destructive",
  },
  503: {
    icon: ServerCrash,
    title: "Service Unavailable",
    description: "The service is temporarily unavailable. We're working to restore it as quickly as possible.",
    iconColor: "text-warning",
    iconBg: "from-warning/20 to-warning/5",
    accentColor: "warning",
  },
};

const defaultConfig: ErrorConfig = {
  icon: AlertTriangle,
  title: "Something Went Wrong",
  description: "An unexpected error occurred. Please try again or contact support if the problem persists.",
  iconColor: "text-destructive",
  iconBg: "from-destructive/20 to-destructive/5",
  accentColor: "destructive",
};

interface AxiosError extends Error {
  response?: {
    status?: number;
    data?: unknown;
  };
  status?: number;
  statusCode?: number;
}

interface ErrorPageProps {
  error?: AxiosError;
  status?: number;
  title?: string;
  description?: string;
  showRetry?: boolean;
  onRetry?: () => void;
}

function extractStatusCode(error?: AxiosError, fallbackStatus?: number): number {
  // Direct status prop
  if (fallbackStatus) return fallbackStatus;

  if (error) {
    // Axios error: error.response.status
    if (error.response?.status) return error.response.status;

    // Direct status on error object
    if (error.status) return error.status;
    if (error.statusCode) return error.statusCode;

    // Try to parse from error message (e.g., "Request failed with status code 400")
    const match = error.message?.match(/status code (\d{3})/i);
    if (match) return parseInt(match[1], 10);
  }

  return 500; // Default fallback
}

export function ErrorPage({
  error,
  status,
  title,
  description,
  showRetry = true,
  onRetry,
}: ErrorPageProps) {
  const router = useRouter();

  const statusCode = extractStatusCode(error, status);
  const config = errorConfigs[statusCode] ?? defaultConfig;

  const displayTitle = title ?? config.title;
  const displayDescription = description ?? error?.message ?? config.description;

  const Icon = config.icon;

  const gridPattern = useMemo(() => {
    const lines = [];
    for (let i = 0; i < 20; i++) {
      lines.push(
        <div
          key={`h-${i}`}
          className="absolute left-0 right-0 h-px bg-gradient-to-r from-transparent via-border/30 to-transparent"
          style={{ top: `${i * 5}%` }}
        />,
        <div
          key={`v-${i}`}
          className="absolute top-0 bottom-0 w-px bg-gradient-to-b from-transparent via-border/30 to-transparent"
          style={{ left: `${i * 5}%` }}
        />
      );
    }
    return lines;
  }, []);

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 relative overflow-hidden">
      {/* Background grid pattern */}
      <div className="absolute inset-0 opacity-30 pointer-events-none">
        {gridPattern}
      </div>

      {/* Radial gradient overlay */}
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,transparent_0%,hsl(var(--background))_70%)] pointer-events-none" />

      <div className="w-full max-w-lg relative z-10">
        {/* Error Code Display */}
        <div className="text-center mb-8">
          <div className="relative inline-block">
            {/* Glowing backdrop for status code */}
            <div
              className={`absolute inset-0 blur-3xl opacity-20 bg-${config.accentColor}`}
              style={{ transform: "scale(2)" }}
            />
            <span className="relative text-[8rem] font-mono font-bold leading-none tracking-tighter text-foreground/10 select-none">
              {statusCode}
            </span>
          </div>
        </div>

        {/* Main Card */}
        <Card className="relative overflow-hidden border-2 border-border/50 backdrop-blur-sm">
          {/* Decorative gradient */}
          <div className={`absolute inset-0 bg-gradient-to-br ${config.iconBg} opacity-50 pointer-events-none`} />

          {/* Animated scan line effect */}
          <div className="absolute inset-0 overflow-hidden pointer-events-none">
            <div className="absolute inset-x-0 h-[2px] bg-gradient-to-r from-transparent via-primary to-transparent animate-scan-line" />
          </div>

          <CardContent className="relative p-8">
            {/* Icon */}
            <div className="flex justify-center mb-6">
              <div className={`w-20 h-20 rounded-2xl bg-gradient-to-br ${config.iconBg} border border-border/50 flex items-center justify-center`}>
                <Icon className={`w-10 h-10 ${config.iconColor}`} />
              </div>
            </div>

            {/* Title & Description */}
            <div className="text-center mb-8">
              <h1 className="text-2xl font-bold tracking-tight mb-2 font-mono">
                {displayTitle}
              </h1>
              <p className="text-muted-foreground text-sm leading-relaxed max-w-sm mx-auto">
                {displayDescription}
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

              {showRetry && onRetry && (
                <Button
                  variant="outline"
                  onClick={onRetry}
                  className="gap-2"
                >
                  <RefreshCw className="w-4 h-4" />
                  Try Again
                </Button>
              )}

              <Button variant="glow" asChild className="gap-2">
                <Link to="/dashboard">
                  <Home className="w-4 h-4" />
                  Dashboard
                </Link>
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Error details for developers */}
        {error && process.env.NODE_ENV === "development" && (
          <details className="mt-6 text-xs">
            <summary className="text-muted-foreground cursor-pointer hover:text-foreground transition-colors">
              Technical Details
            </summary>
            <pre className="mt-2 p-4 rounded-lg bg-muted/50 border border-border/50 overflow-auto text-muted-foreground font-mono">
              {error.stack ?? error.message}
            </pre>
          </details>
        )}
      </div>
    </div>
  );
}
