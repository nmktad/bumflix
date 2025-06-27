export interface Film {
  id: number;
  title: string;
  director: string;
  year: number;
  duration: number;
  genre: string;
  poster: string;
  description: string;
  featured?: boolean;
  rating?: number;
}
