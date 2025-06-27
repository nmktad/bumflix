import FilmCard from "@/components/custom/film-card";
import { collections, featuredFilms, newReleases } from "./data";

export default function Home() {
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          {[...featuredFilms, ...newReleases, ...collections].map((film) => (
            <FilmCard key={film.id} film={film} />
          ))}
        </div>
      </main>
    </div>
  );
}
