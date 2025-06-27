import { collections, featuredFilms, newReleases } from "@/app/data";
import Player from "./player";

export default async function FilmViewPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const selectedFilm = [
    ...featuredFilms,
    ...newReleases,
    ...collections,
  ].filter((film) => film.id == parseInt(id))[0];

  return <Player film={selectedFilm} />;
}
