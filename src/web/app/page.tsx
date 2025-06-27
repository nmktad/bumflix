import FilmCard from "@/components/custom/film-card";
import { collections, featuredFilms, newReleases } from "./data";
import Link from "next/link";

export type Film = {
  title: string;
  slug: string;
  poster_url?: string;
};

async function getMovies() {
  try {
    const res = await fetch("http://localhost:8080/films"); // your Go backend URL
    if (!res.ok) throw new Error("Failed to fetch films");

    const data = await res.json();

    return data.videos ?? [];
  } catch (error) {
    console.error("SSR fetch error:", error);
    return {
      props: {
        films: [],
      },
    };
  }
}

export default async function Home() {
  const movies = (await getMovies()) as Film[];

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          {/* {[...featuredFilms, ...newReleases, ...collections].map((film) => ( */}
          {/*   <FilmCard key={film.id} film={film} /> */}
          {/* ))} */}
          {movies.map((movie) => (
            <Link key={movie.slug} href={`film/${movie.slug}`}>
              {movie.title}
            </Link>
          ))}
        </div>
      </main>
    </div>
  );
}
