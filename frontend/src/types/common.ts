export interface CanvasCardData {
    id: string;
    title: string;
    stars: number;
    owner: {
        id: string;
        username: string;
        avatar: string;
        url: string;
    }
    description: string;
    image: string;
    createdAt: Date;
    updatedAt: Date;
}
