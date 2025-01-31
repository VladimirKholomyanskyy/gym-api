import { useState } from "react";
import { Input, Select, Stack, Flex, Box } from "@chakra-ui/react";
import {
  SelectContent,
  SelectItem,
  SelectRoot,
  SelectTrigger,
  SelectValueText,
} from "../ui/select";

interface SearchAndFilterProps {
  onSearch: (query: string) => void;
  onFilter: (muscle: string, equipment: string) => void;
  muscles: string[];
  equipment: string[];
}

const SearchAndFilterExercises = ({
  onSearch,
  onFilter,
  muscles,
  equipment,
}: SearchAndFilterProps) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedMuscle, setSelectedMuscle] = useState("");
  const [selectedEquipment, setSelectedEquipment] = useState("");

  const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    setSearchQuery(value);
    onSearch(value);
  };

  const handleMuscleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const value = event.target.value;
    setSelectedMuscle(value);
    onFilter(value, selectedEquipment);
  };

  const handleEquipmentChange = (
    event: React.ChangeEvent<HTMLSelectElement>
  ) => {
    const value = event.target.value;
    setSelectedEquipment(value);
    onFilter(selectedMuscle, value);
  };

  return (
    <Box width="100%">
      <Flex direction={{ base: "column", md: "row" }} gap={4} align="center">
        {/* Search Input */}
        <Input
          placeholder="Search exercises..."
          value={searchQuery}
          onChange={handleSearch}
          bg="gray.700"
          color="white"
          _placeholder={{ color: "gray.400" }}
        />

        {/* Filter: Primary Muscle */}
        <SelectRoot
          placeholder="Filter by muscle"
          onValueChange={handleMuscleChange}
          value={selectedMuscle}
          onChange={handleMuscleChange}
          variant="filled"
          bg="gray.700"
          color="white"
        >
          {muscles.map((muscle) => (
            <option key={muscle} value={muscle}>
              {muscle}
            </option>
          ))}
          <SelectTrigger clearable>
            <SelectValueText placeholder="Select Exercise" />
          </SelectTrigger>
          <SelectContent>
            {muscles.map((muscle) => (
              <SelectItem key={muscle} item={muscle}>
                {muscle}
              </SelectItem>
            ))}
          </SelectContent>
        </SelectRoot>

        {/* Filter: Equipment */}
        <SelectRoot
          placeholder="Filter by equipment"
          value={selectedEquipment}
          onChange={handleEquipmentChange}
          variant="filled"
          bg="gray.700"
          color="white"
        >
          {equipment.map((item) => (
            <option key={item} value={item}>
              {item}
            </option>
          ))}
        </SelectRoot>
      </Flex>
    </Box>
  );
};

export default SearchAndFilterExercises;
