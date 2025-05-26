export default function PageHeader({ title, description }: { title: string, description: string}) {
    return (
        <header className="mb-5">
            <h1 className="font-bold text-4xl mb-4">{title}</h1>
            <p className="text-muted-foreground">{description}</p>
        </header>
    );
}
