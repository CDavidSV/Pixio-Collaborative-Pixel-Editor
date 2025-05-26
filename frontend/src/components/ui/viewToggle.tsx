import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { LayoutGrid, StretchHorizontal } from "lucide-react";

export function ViewToggle() {
    return (
        <ToggleGroup type='single' defaultValue='card'>
            <ToggleGroupItem value='table'>
                <StretchHorizontal />
            </ToggleGroupItem>
            <ToggleGroupItem value='card'>
                <LayoutGrid />
            </ToggleGroupItem>
        </ToggleGroup>
    );
}
