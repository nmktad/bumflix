import Link from "next/link";
import { Clock, Heart, Play, Star } from "lucide-react";

import { Film } from "@/types/types";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export default function FilmCard({ film }: { film: Film }) {
  return (
    <Link
      key={film.id}
      href={`/film/${film.id}`}
      className="group cursor-pointer"
    >
      <div className="bg-white rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100 hover:border-gray-200">
        {/* Image Container */}
        <div className="relative aspect-[3/4] overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200">
          <img
            src={film.poster || "/placeholder.svg?height=400&width=300"}
            alt={film.title}
            className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
          />

          {/* Overlay */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />

          {/* Play Button */}
          <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300">
            <Button
              size="lg"
              className="bg-white/90 hover:bg-white text-gray-900 rounded-full shadow-lg backdrop-blur-sm"
            >
              <Play className="w-5 h-5 fill-current" />
            </Button>
          </div>

          {/* Favorite Button */}
          <Button
            variant="ghost"
            size="icon"
            className="absolute top-3 right-3 bg-white/80 hover:bg-white text-gray-700 rounded-full shadow-sm opacity-0 group-hover:opacity-100 transition-opacity duration-300"
          >
            <Heart
              className={cn(
                "w-4 h-4",
                // favorites.has(film.id) && "fill-pink-500 text-pink-500",
              )}
            />
          </Button>

          {/* Duration Badge */}
          <div className="absolute bottom-4 left-4 right-4 space-y-2 opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-y-4 group-hover:translate-y-0">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-1 backdrop-blur-xl bg-black/40 border border-white/20 px-2 py-1 rounded-full">
                <Clock className="w-3 h-3 text-white" />
                <span className="text-white text-xs font-semibold">
                  {film.duration}m
                </span>
              </div>
              {film.rating && (
                <div className="flex items-center space-x-1 backdrop-blur-xl bg-yellow-500/20 border border-yellow-400/30 px-3 py-1 rounded-full">
                  <Star className="w-3 h-3 text-yellow-400 fill-current" />
                  <span className="text-yellow-300 text-xs font-semibold">
                    {film.rating}
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Featured Badge */}
          {film.featured && (
            <div className="absolute top-3 left-3 bg-gradient-to-r from-yellow-400 to-orange-500 text-white px-3 py-1 rounded-full text-xs font-bold">
              Featured
            </div>
          )}
        </div>

        {/* Content */}
        <div className="p-6 space-y-3">
          <div className="space-y-2">
            <h3 className="font-bold text-lg text-gray-900 group-hover:text-gray-700 transition-colors line-clamp-1">
              {film.title}
            </h3>
            <p className="text-gray-600 text-sm">
              {film.director} • {film.year}
            </p>
          </div>

          {film.rating && (
            <div className="flex items-center space-x-1">
              <Star className="w-4 h-4 text-yellow-400 fill-current" />
              <span className="text-sm font-medium text-gray-700">
                {film.rating}
              </span>
            </div>
          )}

          <div className="inline-block bg-gradient-to-r from-pink-100 to-purple-100 text-pink-700 px-3 py-1 rounded-full text-xs font-medium">
            {film.genre}
          </div>

          <p className="text-sm text-gray-500 line-clamp-2 leading-relaxed">
            {film.description}
          </p>
        </div>
      </div>
    </Link>
  );
}
