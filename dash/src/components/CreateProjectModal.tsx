import { useEffect, useState } from 'react';
import { toast } from 'sonner';

interface CreateProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: { name: string; description: string; tags: string[] }) => void;
}

export const CreateProjectModal = ({ isOpen, onClose, onSubmit }: CreateProjectModalProps) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    tagInput: '',
    tags: [] as string[]
  });

  useEffect(() => {
    if (!isOpen) {
      setFormData({
        name: '',
        description: '',
        tagInput: '',
        tags: []
      });
    }
  }, [isOpen]);

  const handleAddTag = () => {
    if (formData.tagInput.trim()) {
      if (formData.tags.includes(formData.tagInput.trim())) {
        toast.error('Tag already exists');
        return;
      }
      setFormData(prev => ({
        ...prev,
        tags: [...prev.tags, prev.tagInput.trim()],
        tagInput: ''
      }));
    }
  };

  const handleRemoveTag = (tagToRemove: string) => {
    setFormData(prev => ({
      ...prev,
      tags: prev.tags.filter(tag => tag !== tagToRemove)
    }));
  };



  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-[#161B22] border border-[#30363D] rounded-lg p-6 w-full max-w-md">
        <h2 className="text-[#C9D1D9] text-xl font-semibold mb-4">Create New Project</h2>

        <form onSubmit={(e) => {
          e.preventDefault();
          onSubmit({
            name: formData.name,
            description: formData.description,
            tags: formData.tags
          });
        }}>
          <div className="space-y-4">
            <div>
              <label className="block text-[#8B949E] text-sm mb-2">Project Name</label>
              <input
                type="text"
                className="w-full px-3 py-2 bg-[#0D1117] border border-[#30363D] rounded-md text-[#C9D1D9] focus:outline-none focus:border-[#1F6FEB]"
                value={formData.name}
                onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                required
              />
            </div>

            <div>
              <label className="block text-[#8B949E] text-sm mb-2">Description</label>
              <textarea
                className="w-full px-3 py-2 bg-[#0D1117] border border-[#30363D] rounded-md text-[#C9D1D9] focus:outline-none focus:border-[#1F6FEB] h-32 resize-none"
                value={formData.description}
                onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
              />
            </div>

            <div>
              <label className="block text-[#8B949E] text-sm mb-2">Tags</label>
              <div className="flex gap-2 mb-2">
                <input
                  type="text"
                  className="flex-1 px-3 py-2 bg-[#0D1117] border border-[#30363D] rounded-md text-[#C9D1D9] focus:outline-none focus:border-[#1F6FEB]"
                  value={formData.tagInput}
                  onChange={(e) => setFormData(prev => ({ ...prev, tagInput: e.target.value }))}
                  onKeyPress={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      handleAddTag();
                    }
                  }}
                  placeholder="Add a tag"
                />
                <button
                  type="button"
                  onClick={handleAddTag}
                  className="px-4 py-2 bg-[#21262D] text-[#C9D1D9] rounded-md hover:bg-[#30363D]"
                >
                  Add
                </button>
              </div>
              <div className="flex flex-wrap gap-2">
                {formData.tags.map(tag => (
                  <span
                    key={tag}
                    className="px-2 py-1 text-sm rounded-full bg-[#1F6FEB33] text-[#1F6FEB] flex items-center gap-2"
                  >
                    {tag}
                    <button
                      type="button"
                      onClick={() => handleRemoveTag(tag)}
                      className="hover:text-[#F85149]"
                    >
                      ×
                    </button>
                  </span>
                ))}
              </div>
            </div>
          </div>

          <div className="flex justify-end gap-3 mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-[#C9D1D9] hover:bg-[#21262D] rounded-md transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-[#1F6FEB] text-white rounded-md hover:bg-[#1A73E8] transition-colors"
            >
              Create Project
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
